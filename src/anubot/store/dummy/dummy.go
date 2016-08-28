package dummy

import (
	"sync"

	"github.com/satori/go.uuid"

	"anubot/store"
	"anubot/twitch/oauth"
)

type Dummy struct {
	mu     sync.Mutex
	users  users
	nonces map[string]struct{}
}

type users map[string]credentials

func (u users) lookup(username string) (credentials, bool) {
	for id, creds := range u {
		if creds.username == username {
			return creds, true
		}
	}
	return credentials{}, false
}

func (u users) exists(username string) bool {
	_, exists := u.lookup(username)
	return exists
}

type credentials struct {
	username string
	password string
}

func New() *Dummy {
	return &Dummy{
		users:  make(users),
		nonces: make(map[string]struct{}),
	}
}

// RegisterUser registers a new user.
func (d *Dummy) RegisterUser(username, password string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.users.exists(username) {
		return "", store.UsernameTaken
	}

	id := uuid.NewV4().String()
	d.users[id] = credentials{
		username: username,
		password: password,
	}
	return id, nil
}

// AuthenticateUser checks to see if the given user credentials are valid.
func (d *Dummy) AuthenticateUser(username, password string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	c, exists := d.users.lookup(username)
	if !exists {
		return flase
	}
	return c.password == password
}

// CreateOauthNonce creates and returns a unique oauth nonce.
func (d *Dummy) CreateOauthNonce() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	nonce := oauth.GenerateNonce()
	d.nonces[nonce] = struct{}{}
	return nonce
}

// OauthNonceExists tells you if the provided nonce was recently created by
// this server.
func (d *Dummy) OauthNonceExists(nonce string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.nonces[nonce]
	return ok
}
