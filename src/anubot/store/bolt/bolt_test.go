package bolt

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/a8m/expect"

	"anubot/store"
	"anubot/twitch"
)

func TestRegisteringAUserReservesThatUsername(t *testing.T) {
	expect := expect.New(t)
	b, cleanup := setupDB(t)
	defer cleanup()

	_, err := b.RegisterUser("test-user", "test-pass")
	expect(err).To.Be.Nil().Else.FailNow()
	_, err = b.RegisterUser("test-user", "test-pass")
	expect(err).To.Equal(store.UsernameTaken)
}

func TestAuthenticationWorks(t *testing.T) {
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

func setupDB(t *testing.T) (*Bolt, func()) {
	path, tmpFileCleanup := tempFile(t)
	b, err := New(path, twitch.API{})
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
