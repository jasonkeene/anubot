package bolt

import (
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"

	"anubot/store"
	"anubot/twitch/oauth"
)

// Bolt is a store backend for boltdb.
type Bolt struct {
	db *bolt.DB
}

// New creates a new bolt store.
func New(path string) (*Bolt, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: time.Second,
	})
	if err != nil {
		return nil, err
	}
	b := &Bolt{
		db: db,
	}
	err = b.createBuckets()
	if err != nil {
		db.Close()
		return nil, err
	}
	return b, nil
}

func (b *Bolt) createBuckets() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("nonces"))
		if err != nil {
			return err
		}
		return nil
	})
}

// Close cleans up the boltdb resources.
func (b *Bolt) Close() error {
	return b.db.Close()
}

// RegisterUser registers a new user returning the user ID.
func (b *Bolt) RegisterUser(username, password string) (string, error) {
	userID := uuid.NewV4().String()
	ur := userRecord{
		UserID:   userID,
		Username: username,
		Password: password,
	}

	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := getUserRecordByUsername(username, tx)
		if err == nil {
			return store.ErrUsernameTaken
		}
		return upsertUserRecord(ur, tx)
	})
	if err != nil {
		return "", err
	}

	return userID, nil
}

// AuthenticateUser checks to see if the given user credentials are valid. If
// they are the user ID is returned with a bool to indicate success.
func (b *Bolt) AuthenticateUser(username, password string) (string, bool) {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecordByUsername(username, tx)
		return err
	})
	if err != nil {
		return "", false
	}

	if ur.Password != password {
		return "", false
	}
	return ur.UserID, true
}

// CreateOauthNonce creates and returns a unique oauth nonce.
func (b *Bolt) CreateOauthNonce(userID string, tu store.TwitchUser) (string, error) {
	switch tu {
	case store.Streamer:
	case store.Bot:
	default:
		return "", store.ErrInvalidTwitchUserType
	}

	nonce := oauth.GenerateNonce()
	nr := nonceRecord{
		Nonce:   nonce,
		UserID:  userID,
		TU:      tu,
		Created: time.Now(),
	}

	err := b.db.Update(func(tx *bolt.Tx) error {
		return upsertNonceRecord(nr, tx)
	})
	if err != nil {
		return "", err
	}

	return nonce, nil
}

// OauthNonceExists tells you if the provided nonce was recently created and
// not yet finished.
func (b *Bolt) OauthNonceExists(nonce string) bool {
	err := b.db.View(func(tx *bolt.Tx) error {
		_, err := getNonceRecord(nonce, tx)
		return err
	})
	if err != nil {
		return false
	}
	return true
}

// FinishOauthNonce completes the oauth flow, removing the nonce and storing
// the oauth data.
func (b *Bolt) FinishOauthNonce(nonce, username string, od oauth.Data) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		nr, err := getNonceRecord(nonce, tx)
		if err != nil {
			return err
		}

		ur, err := getUserRecord(nr.UserID, tx)
		if err != nil {
			return err
		}

		switch nr.TU {
		case store.Streamer:
			ur.StreamerOD = od
			ur.StreamerUsername = username
		case store.Bot:
			ur.BotOD = od
			ur.BotUsername = username
		default:
			return errors.New("bad twitch user type, this should never happen")
		}

		err = deleteNonceRecord(nr.Nonce, tx)
		if err != nil {
			return err
		}

		return upsertUserRecord(ur, tx)
	})
}

// TwitchStreamerAuthenticated tells you if the user has authenticated with
// twitch and that we have valid oauth credentials.
func (b *Bolt) TwitchStreamerAuthenticated(userID string) bool {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return false
	}

	return ur.StreamerOD.AccessToken != ""
}

// TwitchStreamerCredentials gives you the credentials for the streamer user.
func (b *Bolt) TwitchStreamerCredentials(userID string) (string, string) {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return "", ""
	}

	return ur.StreamerUsername, ur.StreamerOD.AccessToken
}

// TwitchBotAuthenticated tells you if the user has authenticated his bot with
// twitch and that we have valid oauth credentials.
func (b *Bolt) TwitchBotAuthenticated(userID string) bool {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return false
	}

	return ur.BotOD.AccessToken != ""
}

// TwitchBotCredentials gives you the credentials for the streamer user.
func (b *Bolt) TwitchBotCredentials(userID string) (string, string) {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return "", ""
	}

	return ur.BotUsername, ur.BotOD.AccessToken
}

// TwitchAuthenticated tells you if the user has authenticated his bot and
// his streamer user with twitch and that we have valid oauth credentials.
func (b *Bolt) TwitchAuthenticated(userID string) bool {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return false
	}

	return ur.StreamerOD.AccessToken != "" && ur.BotOD.AccessToken != ""
}

// TwitchClearAuth removes all the auth data for twitch for the user.
func (b *Bolt) TwitchClearAuth(userID string) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		ur, err := getUserRecord(userID, tx)
		if err != nil {
			return err
		}
		ur.StreamerUsername = ""
		ur.StreamerOD = oauth.Data{}
		ur.BotUsername = ""
		ur.BotOD = oauth.Data{}
		return upsertUserRecord(ur, tx)
	})
	if err != nil {
		log.Printf("could not clear twitch auth: %s", err)
	}
}
