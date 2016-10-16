package bolt

import (
	"encoding/json"
	"errors"
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
func (b *Bolt) RegisterUser(username, password string) (userID string, err error) {
	userID = uuid.NewV4().String()
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		existingData := b.Get([]byte(username))
		if existingData != nil {
			return store.UsernameTaken
		}

		ur := userRecord{
			UserID:   userID,
			Username: username,
			Password: password,
		}
		urb, err := json.Marshal(ur)
		if err != nil {
			return err
		}

		return b.Put([]byte(username), urb)
	})

	if err != nil {
		return "", err
	}
	return userID, nil
}

// AuthenticateUser checks to see if the given user credentials are valid.
func (b *Bolt) AuthenticateUser(username, password string) (userID string, authenticated bool) {
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		data := b.Get([]byte(username))
		var ur userRecord
		err := json.Unmarshal(data, &ur)
		if err != nil {
			return err
		}
		if ur.Password != password {
			return errors.New("bad auth")
		}
		userID = ur.UserID
		return nil
	})

	if err != nil {
		return "", false
	}
	return userID, true
}

// CreateOauthNonce creates and returns a unique oauth nonce.
func (b *Bolt) CreateOauthNonce(userID string, tu store.TwitchUser) (nonce string) {
	// TODO: implement
	return ""
}

// OauthNonceExists tells you if the provided nonce was recently created and
// not yet finished.
func (b *Bolt) OauthNonceExists(nonce string) (exists bool) {
	// TODO: implement
	return false
}

// FinishOauthNonce completes the oauth flow, removing the nonce and storing
// the oauth data.
func (b *Bolt) FinishOauthNonce(nonce, username string, od oauth.Data) error {
	// TODO: implement
	return nil
}

// TwitchStreamerAuthenticated tells you if the user has authenticated with
// twitch and that we have valid oauth credentials.
func (b *Bolt) TwitchStreamerAuthenticated(userID string) bool {
	// TODO: implement
	return false
}

// TwitchStreamerCredentials gives you the credentials for the streamer user.
func (b *Bolt) TwitchStreamerCredentials(userID string) (string, string) {
	// TODO: implement
	return "", ""
}

// TwitchBotAuthenticated tells you if the user has authenticated his bot with
// twitch and that we have valid oauth credentials.
func (b *Bolt) TwitchBotAuthenticated(userID string) bool {
	// TODO: implement
	return false
}

// TwitchBotCredentials gives you the credentials for the streamer user.
func (b *Bolt) TwitchBotCredentials(userID string) (string, string) {
	// TODO: implement
	return "", ""
}

type userRecord struct {
	UserID           string     `json:"user_id"`
	Username         string     `json:"username"`
	Password         string     `json:"password"`
	StreamerUsername string     `json:"streamer_username"`
	StreamerOD       oauth.Data `json:"streamer_od"`
	BotUsername      string     `json:"bot_username"`
	BotOD            oauth.Data `json:"bot_od"`
}

type nonceRecord struct {
	nonce   string
	userID  string
	tu      store.TwitchUser
	created time.Time
}
