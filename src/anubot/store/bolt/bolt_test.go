package bolt

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/a8m/expect"
	"github.com/fluffle/goirc/client"

	"anubot/store"
	"anubot/stream"
)

func TestThatBoltBackendCompliesWithAllStoreMethods(t *testing.T) {
	var _ store.Store = &Bolt{}
}

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

	od := store.OauthData{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	err = b.FinishOauthNonce(nonce, "test-streamer-user", 12345, od)
	expect(err).To.Be.Nil()
	expect(b.OauthNonceExists(nonce)).Not.To.Be.Ok()
	expect(b.TwitchStreamerAuthenticated(userID)).To.Be.Ok()

	user, pass, _ := b.TwitchStreamerCredentials(userID)
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

	od := store.OauthData{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	err = b.FinishOauthNonce(nonce, "test-bot-user", 12345, od)
	expect(err).To.Be.Nil()
	expect(b.OauthNonceExists(nonce)).Not.To.Be.Ok()
	expect(b.TwitchBotAuthenticated(userID)).To.Be.Ok()

	user, pass, _ := b.TwitchBotCredentials(userID)
	expect(user).To.Equal("test-bot-user")
	expect(pass).To.Equal("test-access-token")
}

func TestThatYouCanClearTwitchAuthData(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	userID, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil()

	od := store.OauthData{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	nonce, err := b.CreateOauthNonce(userID, store.Streamer)
	expect(err).To.Be.Nil()
	err = b.FinishOauthNonce(nonce, "test-streamer-user", 12345, od)
	expect(err).To.Be.Nil()
	expect(b.TwitchStreamerAuthenticated(userID)).To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).Not.To.Be.Ok()

	nonce, err = b.CreateOauthNonce(userID, store.Bot)
	expect(err).To.Be.Nil()
	err = b.FinishOauthNonce(nonce, "test-bot-user", 12345, od)
	expect(err).To.Be.Nil()
	expect(b.TwitchBotAuthenticated(userID)).To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).To.Be.Ok()

	b.TwitchClearAuth(userID)
	expect(b.TwitchStreamerAuthenticated(userID)).Not.To.Be.Ok()
	expect(b.TwitchBotAuthenticated(userID)).Not.To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).Not.To.Be.Ok()
}

func TestThatYouCanStoreMessages(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	userID, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil()

	od := store.OauthData{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scope:        []string{"test-scope"},
	}
	nonce, err := b.CreateOauthNonce(userID, store.Streamer)
	expect(err).To.Be.Nil()
	err = b.FinishOauthNonce(nonce, "test-streamer-user", 12345, od)
	expect(err).To.Be.Nil()
	expect(b.TwitchStreamerAuthenticated(userID)).To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).Not.To.Be.Ok()

	nonce, err = b.CreateOauthNonce(userID, store.Bot)
	expect(err).To.Be.Nil()
	err = b.FinishOauthNonce(nonce, "test-bot-user", 54321, od)
	expect(err).To.Be.Nil()
	expect(b.TwitchBotAuthenticated(userID)).To.Be.Ok()
	expect(b.TwitchAuthenticated(userID)).To.Be.Ok()

	msg1 := stream.RXMessage{
		Type: stream.Twitch,
		Twitch: &stream.RXTwitch{
			OwnerID: 12345,
			Line: &client.Line{
				Raw: "test-message-1",
			},
		},
	}
	err = b.StoreMessage(msg1)
	expect(err).To.Be.Nil().Else.FailNow()
	msg2 := stream.RXMessage{
		Type: stream.Twitch,
		Twitch: &stream.RXTwitch{
			OwnerID: 12345,
			Line: &client.Line{
				Raw: "test-message-2",
			},
		},
	}
	err = b.StoreMessage(msg2)
	expect(err).To.Be.Nil().Else.FailNow()

	messages, err := b.FetchRecentMessages(userID)
	expect(err).To.Be.Nil().Else.FailNow()
	expect(len(messages)).To.Equal(2).Else.FailNow()
	expect(messages[0].Twitch.Line.Raw).To.Equal("test-message-1")
	expect(messages[1].Twitch.Line.Raw).To.Equal("test-message-2")
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
