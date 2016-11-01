package bolt

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"

	"anubot/store"
	"anubot/stream"
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
		closeErr := db.Close()
		if closeErr != nil {
			log.Printf("got an error closing bolt db while in error state: %s", err)
		}
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
		_, err = tx.CreateBucketIfNotExists([]byte("messages"))
		return err
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
	return err == nil
}

// FinishOauthNonce completes the oauth flow, removing the nonce and storing
// the oauth data.
func (b *Bolt) FinishOauthNonce(
	nonce string,
	username string,
	userID int,
	od store.OauthData,
) error {
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
			ur.StreamerID = userID
		case store.Bot:
			ur.BotOD = od
			ur.BotUsername = username
			ur.BotID = userID
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
func (b *Bolt) TwitchStreamerCredentials(userID string) (string, string, int) {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return "", "", 0
	}

	return ur.StreamerUsername, ur.StreamerOD.AccessToken, ur.StreamerID
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
func (b *Bolt) TwitchBotCredentials(userID string) (string, string, int) {
	var ur userRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		ur, err = getUserRecord(userID, tx)
		return err
	})
	if err != nil {
		return "", "", 0
	}

	return ur.BotUsername, ur.BotOD.AccessToken, ur.BotID
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
		ur.StreamerOD = store.OauthData{}
		ur.BotUsername = ""
		ur.BotOD = store.OauthData{}
		return upsertUserRecord(ur, tx)
	})
	if err != nil {
		log.Printf("could not clear twitch auth: %s", err)
	}
}

// StoreMessage stores a message for a given user for later searching and
// scrollback history.
func (b *Bolt) StoreMessage(msg stream.RXMessage) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return upsertMessage(msg, tx)
	})
}

type byTimestamp []stream.RXMessage

func (a byTimestamp) Len() int           { return len(a) }
func (a byTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTimestamp) Less(i, j int) bool { return a[i].Twitch.Line.Time.Before(a[j].Twitch.Line.Time) }

// FetchRecentMessages gets the recent messages for the user's channel.
func (b *Bolt) FetchRecentMessages(userID string) ([]stream.RXMessage, error) {
	if !b.TwitchAuthenticated(userID) {
		return nil, errors.New("user must be authenticated via twitch before requesting recent messages")
	}

	var messages []stream.RXMessage
	_, _, streamerUserID := b.TwitchStreamerCredentials(userID)
	_, _, botUserID := b.TwitchBotCredentials(userID)

	var mr messageRecord
	err := b.db.View(func(tx *bolt.Tx) error {
		var err error
		mr, err = getMessageRecord("twitch:"+strconv.Itoa(streamerUserID), tx)
		return err
	})
	if err != nil {
		log.Printf("could not query messages for streamer: %s", err)
	}

	messages = []stream.RXMessage(mr)

	err = b.db.View(func(tx *bolt.Tx) error {
		var err error
		mr, err = getMessageRecord("twitch:"+strconv.Itoa(botUserID), tx)
		return err
	})
	if err != nil {
		log.Printf("could not query messages for bot: %s", err)
	}

	messages = append(messages, []stream.RXMessage(mr)...)
	sort.Sort(byTimestamp(messages))
	return messages[:min(len(messages), 500)], nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// QueryMessages allows the user to search for messages that match a
// search string.
func (b *Bolt) QueryMessages(userID, search string) ([]stream.RXMessage, error) {
	panic("not implemented")
}
