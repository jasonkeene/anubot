package stream

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectingOverTLS(t *testing.T) {
	require := require.New(t)
	server := newFakeIRCServer(t)
	defer server.close()
	defer patchTwitch(server.port())()
	d := newMockDispatcher()

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		clientConn, err := connectTwitch("test-user", "test-pass", "#test-chan", d)
		defer clientConn.close()
		if err != nil {
			log.Panic("unable to connect to twitch")
		}
	}()

	serverConn := server.accept()

	pass := serverConn.receive("PASS")
	require.Equal("PASS test-pass", pass)
	nick := serverConn.receive("NICK")
	require.Equal("NICK test-user", nick)
	user := serverConn.receive("USER")
	require.Equal("USER anubot 12 * :test-user", user)

	serverConn.send(":127.0.0.1 001 test-user :GLHF!")

	join := serverConn.receive("JOIN")
	require.Equal("JOIN #test-chan", join)

	serverConn.receive("QUIT")
	serverConn.close()

	<-clientDone
}

func TestDispatchingMessages(t *testing.T) {
	require := require.New(t)
	server := newFakeIRCServer(t)
	defer server.close()
	defer patchTwitch(server.port())()
	d := newMockDispatcher()

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		// racey
		clientConn, err := connectTwitch("test-user", "test-pass", "#test-chan", d)
		defer clientConn.close()
		if err != nil {
			log.Panic("unable to connect to twitch")
		}
	}()

	serverConn, cleanup := acceptConn(server)

	serverConn.send("PRIVMSG #test-chan :test-message")
	msg := <-d.DispatchInput.Message
	require.Equal(msg.Type, Twitch)
	require.Equal(msg.Twitch.Line.Raw, "PRIVMSG #test-chan :test-message")

	cleanup()

	<-clientDone
}

func TestSendingMessages(t *testing.T) {
	require := require.New(t)
	server := newFakeIRCServer(t)
	defer server.close()
	defer patchTwitch(server.port())()
	d := newMockDispatcher()

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		clientConn, err := connectTwitch("test-user", "test-pass", "#test-chan", d)
		defer clientConn.close()
		if err != nil {
			log.Panic("unable to connect to twitch")
		}
		clientConn.send(TXMessage{
			To:      "#test-chan",
			Message: "test-message",
		})
	}()

	serverConn, cleanup := acceptConn(server)

	msg := serverConn.receive("PRIVMSG")
	require.Equal(msg, "PRIVMSG #test-chan :test-message")

	cleanup()

	<-clientDone
}

func TestConnectingToUnresponsiveServer(t *testing.T) {
	require := require.New(t)
	server := newFakeIRCServer(t)
	defer patchTwitch(server.port())()
	d := newMockDispatcher()
	server.close()

	_, err := connectTwitch("test-user", "test-pass", "#test-chan", d)
	require.Error(err)
}

func patchTwitch(port int) func() {
	oHost, oPort := twitchHost, twitchPort
	oSkip, oFlood := insecureSkipVerify, flood
	twitchHost, twitchPort = "127.0.0.1", port
	insecureSkipVerify, flood = true, true
	return func() {
		twitchHost, twitchPort = oHost, oPort
		insecureSkipVerify, flood = oSkip, oFlood
	}
}

func acceptConn(server *fakeIRCServer) (*ircConn, func()) {
	serverConn := server.accept()

	serverConn.receive("PASS")
	serverConn.receive("NICK")
	serverConn.receive("USER")
	serverConn.send(":127.0.0.1 001 test-user :GLHF!")
	serverConn.receive("JOIN")

	return serverConn, func() {
		serverConn.receive("QUIT")
		serverConn.close()
	}
}
