package bolt

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/a8m/expect"

	"anubot/store"
	"anubot/twitch/oauth"
)

func TestThatRegisteringAUserReservesThatUsername(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	_, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil().Else.FailNow()
	_, err = b.RegisterUser("test-user", "test-pass")
	expect(err).To.Equal(store.ErrUsernameTaken)
}

func TestThatUsersCanAuthenticate(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	expectedUserID, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil()

	userID, authenticated := b.AuthenticateUser("test-user", "bad-pass")
	expect(userID).To.Equal("")
	expect(authenticated).Not.To.Be.Ok()

	userID, authenticated = b.AuthenticateUser("test-user", "test-pass")
	expect(userID).To.Equal(expectedUserID)
	expect(authenticated).To.Be.Ok()
}

func TestThatStreamerUsersCanAuthenticate(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	userID, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil()

	nonce, err := b.CreateOauthNonce(userID, store.Streamer)
	expect(err).To.Be.Nil()
	expect(b.OauthNonceExists(nonce)).To.Be.Ok()

	od := oauth.Data{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	err = b.FinishOauthNonce(nonce, "test-streamer-user", od)
	expect(err).To.Be.Nil()
	expect(b.OauthNonceExists(nonce)).Not.To.Be.Ok()
	expect(b.TwitchStreamerAuthenticated(userID)).To.Be.Ok()

	user, pass := b.TwitchStreamerCredentials(userID)
	expect(user).To.Equal("test-streamer-user")
	expect(pass).To.Equal("test-access-token")
}

func TestThatOauthFlowForBotsWorks(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	userID, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil()

	nonce, err := b.CreateOauthNonce(userID, store.Bot)
	expect(err).To.Be.Nil()
	expect(b.OauthNonceExists(nonce)).To.Be.Ok()

	od := oauth.Data{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	err = b.FinishOauthNonce(nonce, "test-bot-user", od)
	expect(err).To.Be.Nil()
	expect(b.OauthNonceExists(nonce)).Not.To.Be.Ok()
	expect(b.TwitchBotAuthenticated(userID)).To.Be.Ok()

	user, pass := b.TwitchBotCredentials(userID)
	expect(user).To.Equal("test-bot-user")
	expect(pass).To.Equal("test-access-token")
}

func TestThatYouCanClearTwitchAuthData(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	userID, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil()

	od := oauth.Data{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	nonce, err := b.CreateOauthNonce(userID, store.Streamer)
	expect(err).To.Be.Nil()
	err = b.FinishOauthNonce(nonce, "test-streamer-user", od)
	expect(err).To.Be.Nil()
	expect(b.TwitchStreamerAuthenticated(userID)).To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).Not.To.Be.Ok()

	nonce, err = b.CreateOauthNonce(userID, store.Bot)
	expect(err).To.Be.Nil()
	err = b.FinishOauthNonce(nonce, "test-bot-user", od)
	expect(err).To.Be.Nil()
	expect(b.TwitchBotAuthenticated(userID)).To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).To.Be.Ok()

	b.TwitchClearAuth(userID)
	expect(b.TwitchStreamerAuthenticated(userID)).Not.To.Be.Ok()
	expect(b.TwitchBotAuthenticated(userID)).Not.To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).Not.To.Be.Ok()
}

func setupDB(t *testing.T) (*Bolt, func()) {
	path, tmpFileCleanup := tempFile(t)
	b, err := New(path)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	return b, func() {
		b.Close()
		tmpFileCleanup()
	}
}

func tempFile(t *testing.T) (string, func()) {
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Println("could not obtain a temporary file")
		t.FailNow()
	}
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}
