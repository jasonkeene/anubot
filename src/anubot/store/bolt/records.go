package bolt

import (
	"anubot/store"
	"anubot/twitch/oauth"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

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
	Nonce   string           `json:"nonce"`
	UserID  string           `json:"user_id"`
	TU      store.TwitchUser `json:"tu"`
	Created time.Time        `json:"created"`
}

func upsertNonceRecord(nr nonceRecord, tx *bolt.Tx) error {
	b := tx.Bucket([]byte("nonces"))

	nrb, err := json.Marshal(nr)
	if err != nil {
		return err
	}
	return b.Put([]byte(nr.Nonce), nrb)
}

func getNonceRecord(nonce string, tx *bolt.Tx) (nonceRecord, error) {
	b := tx.Bucket([]byte("nonces"))

	read := b.Get([]byte(nonce))
	if read == nil {
		return nonceRecord{}, store.ErrUnknownNonce
	}

	var nr nonceRecord
	err := json.Unmarshal(read, &nr)
	if err != nil {
		return nonceRecord{}, err
	}
	return nr, nil
}

func deleteNonceRecord(nonce string, tx *bolt.Tx) error {
	b := tx.Bucket([]byte("nonces"))

	return b.Delete([]byte(nonce))
}

func upsertUserRecord(ur userRecord, tx *bolt.Tx) error {
	b := tx.Bucket([]byte("users"))

	urb, err := json.Marshal(ur)
	if err != nil {
		return err
	}
	return b.Put([]byte(ur.UserID), urb)
}

func getUserRecord(userID string, tx *bolt.Tx) (userRecord, error) {
	b := tx.Bucket([]byte("users"))

	read := b.Get([]byte(userID))
	if read == nil {
		return userRecord{}, store.ErrUnknownUserID
	}

	var ur userRecord
	err := json.Unmarshal(read, &ur)
	if err != nil {
		return userRecord{}, err
	}
	return ur, nil
}

func getUserRecordByUsername(username string, tx *bolt.Tx) (userRecord, error) {
	b := tx.Bucket([]byte("users"))

	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var ur userRecord
		err := json.Unmarshal(v, &ur)
		if err != nil {
			continue
		}
		if username == ur.Username {
			return ur, nil
		}
	}
	return userRecord{}, store.ErrUnknownUsername
}

func deleteUserRecord(userID string, tx *bolt.Tx) error {
	b := tx.Bucket([]byte("users"))
	return b.Delete([]byte(userID))
}
