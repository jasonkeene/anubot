package stream

import (
	"log"
	"testing"

	"github.com/a8m/expect"
)

func TestConnectingOverTLS(t *testing.T) {
	expect := expect.New(t)
	server := newFakeIRCServer(t)
	defer server.close()
	defer patchTwitch(server.port())()
	d := newMockDispatcher()
	twitch := newMockTwitchUserIDFetcher()
	close(twitch.UserIDOutput.UserID)
	close(twitch.UserIDOutput.Err)

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		clientConn, err := connectTwitch("test-user", "test-pass", "#test-chan", d, twitch)
		defer func() {
			err := clientConn.close()
			if err != nil {
				log.Panic("error in tearing down client conn")
			}
		}()
		if err != nil {
			log.Panic("unable to connect to twitch")
		}
	}()

	serverConn := server.accept()

	pass := serverConn.receive("PASS")
	expect(pass).To.Equal("PASS test-pass")
	nick := serverConn.receive("NICK")
	expect(nick).To.Equal("NICK test-user")
	user := serverConn.receive("USER")
	expect(user).To.Equal("USER anubot 12 * :test-user")

	serverConn.send(":127.0.0.1 001 test-user :GLHF!")

	join := serverConn.receive("JOIN")
	expect(join).To.Equal("JOIN #test-chan")

	cap := serverConn.receive("CAP")
	expect(cap).To.Equal("CAP REQ :" + capString)

	serverConn.receive("QUIT")
	serverConn.close()

	<-clientDone
}

func TestDispatchingMessages(t *testing.T) {
	expect := expect.New(t)
	server := newFakeIRCServer(t)
	defer server.close()
	defer patchTwitch(server.port())()
	d := newMockDispatcher()
	twitch := newMockTwitchUserIDFetcher()
	twitch.UserIDOutput.UserID <- 12345
	close(twitch.UserIDOutput.Err)

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		// racey
		clientConn, err := connectTwitch("test-user", "test-pass", "#test-chan", d, twitch)
		defer func() {
			err := clientConn.close()
			if err != nil {
				log.Panic("error in tearing down client conn")
			}
		}()
		if err != nil {
			log.Panic("unable to connect to twitch")
		}
	}()

	serverConn, cleanup := acceptConn(server)

	serverConn.send("PRIVMSG #test-chan :test-message")
	topic := <-d.DispatchInput.Topic
	expect(topic).To.Equal("twitch:test-user")
	msg := <-d.DispatchInput.Message
	expect(msg.Type).To.Equal(Twitch)
	expect(msg.Twitch.OwnerID).To.Equal(12345)
	expect(msg.Twitch.Line.Raw).To.Equal("PRIVMSG #test-chan :test-message")

	cleanup()

	<-clientDone
}

func TestSendingMessages(t *testing.T) {
	expect := expect.New(t)
	server := newFakeIRCServer(t)
	defer server.close()
	defer patchTwitch(server.port())()
	d := newMockDispatcher()
	twitch := newMockTwitchUserIDFetcher()
	close(twitch.UserIDOutput.UserID)
	close(twitch.UserIDOutput.Err)

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		clientConn, err := connectTwitch("test-user", "test-pass", "#test-chan", d, twitch)
		defer func() {
			err := clientConn.close()
			if err != nil {
				log.Panic("error in tearing down client conn")
			}
		}()
		if err != nil {
			log.Panic("unable to connect to twitch")
		}
		clientConn.send(TXMessage{
			Type: Twitch,
			Twitch: &TXTwitch{
				To:      "#test-chan",
				Message: "test-message",
			},
		})
	}()

	serverConn, cleanup := acceptConn(server)

	msg := serverConn.receive("PRIVMSG")
	expect(msg).To.Equal("PRIVMSG #test-chan :test-message")

	cleanup()

	<-clientDone
}

func TestConnectingToUnresponsiveServer(t *testing.T) {
	expect := expect.New(t)
	server := newFakeIRCServer(t)
	defer patchTwitch(server.port())()
	d := newMockDispatcher()
	twitch := newMockTwitchUserIDFetcher()
	close(twitch.UserIDOutput.UserID)
	close(twitch.UserIDOutput.Err)
	server.close()

	_, err := connectTwitch("test-user", "test-pass", "#test-chan", d, twitch)
	expect(err).Not.To.Be.Nil()
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
