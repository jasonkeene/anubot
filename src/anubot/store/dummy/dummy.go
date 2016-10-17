package dummy

import (
	"errors"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"anubot/store"
	"anubot/twitch/oauth"
)

// Dummy is a store backend that stores everything in memory.
type Dummy struct {
	mu     sync.Mutex
	users  users
	nonces map[string]nonceRecord
}

// New creates a new Dummy store.
func New() *Dummy {
	return &Dummy{
		users:  make(users),
		nonces: make(map[string]nonceRecord),
	}
}

// Close is a NOP on the dummy store.
func (d *Dummy) Close() error {
	return nil
}

// RegisterUser registers a new user returning the user ID.
func (d *Dummy) RegisterUser(username, password string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.users.exists(username) {
		return "", store.ErrUsernameTaken
	}

	id := uuid.NewV4().String()
	d.users[id] = userRecord{
		username: username,
		password: password,
	}
	return id, nil
}

// AuthenticateUser checks to see if the given user credentials are valid. If
// they are the user ID is returned with a bool to indicate success.
func (d *Dummy) AuthenticateUser(username, password string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	id, c, exists := d.users.lookup(username)
	if !exists {
		return "", false
	}
	if c.password != password {
		return "", false
	}
	return id, true
}

// CreateOauthNonce creates and returns a unique oauth nonce.
func (d *Dummy) CreateOauthNonce(userID string, tu store.TwitchUser) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch tu {
	case store.Streamer:
	case store.Bot:
	default:
		return "", errors.New("bad twitch user type in CreateOauthNonce")
	}

	nonce := oauth.GenerateNonce()
	d.nonces[nonce] = nonceRecord{
		userID:  userID,
		tu:      tu,
		created: time.Now(),
	}
	return nonce, nil
}

// OauthNonceExists tells you if the provided nonce was recently created and
// not yet finished.
func (d *Dummy) OauthNonceExists(nonce string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.nonces[nonce]
	return ok
}

// FinishOauthNonce completes the oauth flow, removing the nonce and storing
// the oauth data.
func (d *Dummy) FinishOauthNonce(nonce, username string, od oauth.Data) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	nr, ok := d.nonces[nonce]
	if !ok {
		return store.ErrUnknownNonce
	}

	userRecord := d.users[nr.userID]
	switch nr.tu {
	case store.Streamer:
		userRecord.streamerOD = od
		userRecord.streamerUsername = username
	case store.Bot:
		userRecord.botOD = od
		userRecord.botUsername = username
	default:
		return errors.New("bad twitch user type, this should never happen")
	}

	delete(d.nonces, nonce)
	d.users[nr.userID] = userRecord
	return nil
}

// TwitchStreamerAuthenticated tells you if the user has authenticated with
// twitch and that we have valid oauth credentials.
func (d *Dummy) TwitchStreamerAuthenticated(userID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	userRecord := d.users[userID]
	return userRecord.streamerOD.AccessToken != ""
}

// TwitchStreamerCredentials gives you the credentials for the streamer user.
func (d *Dummy) TwitchStreamerCredentials(userID string) (string, string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ur := d.users[userID]
	return ur.streamerUsername, ur.streamerOD.AccessToken
}

// TwitchBotAuthenticated tells you if the user has authenticated his bot with
// twitch and that we have valid oauth credentials.
func (d *Dummy) TwitchBotAuthenticated(userID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	userRecord := d.users[userID]
	return userRecord.botOD.AccessToken != ""
}

// TwitchBotCredentials gives you the credentials for the streamer user.
func (d *Dummy) TwitchBotCredentials(userID string) (string, string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ur := d.users[userID]
	return ur.botUsername, ur.botOD.AccessToken
}

// TwitchAuthenticated tells you if the user has authenticated his bot and
// his streamer user with twitch and that we have valid oauth credentials.
func (d *Dummy) TwitchAuthenticated(userID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	userRecord := d.users[userID]
	return userRecord.streamerOD.AccessToken != "" &&
		userRecord.botOD.AccessToken != ""
}

// TwitchClearAuth removes all the auth data for twitch for the user.
func (d *Dummy) TwitchClearAuth(userID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	userRecord := d.users[userID]
	userRecord.streamerOD = oauth.Data{}
	userRecord.streamerUsername = ""
	userRecord.botOD = oauth.Data{}
	userRecord.botUsername = ""
}
