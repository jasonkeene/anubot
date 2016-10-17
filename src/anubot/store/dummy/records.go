package dummy

import (
	"anubot/store"
	"anubot/twitch/oauth"
	"time"
)

type nonceRecord struct {
	userID  string
	tu      store.TwitchUser
	created time.Time
}

type users map[string]userRecord

func (u users) lookup(username string) (string, userRecord, bool) {
	for id, creds := range u {
		if creds.username == username {
			return id, creds, true
		}
	}
	return "", userRecord{}, false
}

func (u users) exists(username string) bool {
	_, _, exists := u.lookup(username)
	return exists
}

type userRecord struct {
	username         string
	password         string
	streamerUsername string
	streamerOD       oauth.Data
	botUsername      string
	botOD            oauth.Data
}
