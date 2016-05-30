package bot

import (
	"crypto/tls"
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/fluffle/goirc/client"
)

type ConnConfig struct {
	StreamerUsername string
	StreamerPassword string
	BotUsername      string
	BotPassword      string
	Host             string
	Port             int
	TLSConfig        *tls.Config
}

// Bot communicates with the IRC server and has pointers to features.
type Bot struct {
	mu           sync.Mutex
	connected    bool
	streamerConn *client.Conn
	botConn      *client.Conn
}

// Connect establishes connections to the IRC server.
func (b *Bot) Connect(c *ConnConfig) (chan struct{}, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// check to see if we are connected already
	if b.connected {
		return nil, errors.New("This bot is already connected, disconnect first.")
	}

	// fix up TLS config if not provided
	initTLS(c)

	// connect to the server
	streamerConn, streamerDisconnected, err := connectUser(c.StreamerUsername, c.StreamerPassword, c)
	if err != nil {
		return nil, err
	}
	botConn, botDisconnected, err := connectUser(c.BotUsername, c.BotPassword, c)
	if err != nil {
		streamerConn.Quit()
		return nil, err
	}

	b.connected = true
	b.streamerConn = streamerConn
	b.botConn = botConn

	// signal disconnect on either bot or streamer connection
	disconnected := make(chan struct{})
	go func() {
		defer close(disconnected)
		for {
			select {
			case <-streamerDisconnected:
				return
			case <-botDisconnected:
				return
			}
		}
	}()

	return disconnected, nil
}

// Disconnect tears down the connections to the IRC server and resets the state
// of the bot so that you can connect again.
func (b *Bot) Disconnect() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.connected = false

	if b.streamerConn != nil {
		b.streamerConn.Quit()
		b.streamerConn = nil
	}
	if b.botConn != nil {
		b.botConn.Quit()
		b.botConn = nil
	}
}

// Join joins an IRC channel.
// TODO: Implement parting logic
func (b *Bot) Join(channel string) {
	b.streamerConn.Join(channel)
	b.botConn.Join(channel)
}

// Send sents a raw message to the IRC server.
func (b *Bot) Send(user, message string) {
	switch user {
	case "streamer":
		b.streamerConn.Raw(message)
	case "bot":
		b.botConn.Raw(message)
	default:
		panic("Bad user provided for sending message")
	}
}

func initTLS(c *ConnConfig) {
	if c.TLSConfig == nil {
		c.TLSConfig = &tls.Config{
			ServerName: c.Host,
		}
	}
}

func connectUser(username, password string, c *ConnConfig) (*client.Conn, chan struct{}, error) {
	// create client
	cfg := client.NewConfig(username)
	cfg.Me.Name = username
	cfg.Me.Ident = "anubot"
	cfg.Pass = password
	cfg.SSL = true
	cfg.SSLConfig = c.TLSConfig
	cfg.Server = net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	conn := client.Client(cfg)

	// register dis/connect handlers
	connected := make(chan struct{})
	disconnected := make(chan struct{})
	registerSignalHandler(client.CONNECTED, conn, connected)
	registerSignalHandler(client.DISCONNECTED, conn, disconnected)

	if err := conn.Connect(); err != nil {
		return nil, nil, err
	}
	<-connected
	return conn, disconnected, nil
}

func registerSignalHandler(event string, conn *client.Conn, signal chan struct{}) {
	// TODO: get remover working to prevent closing signal chan multiple times
	//var remover client.Remover
	conn.HandleFunc(event, func(conn *client.Conn, line *client.Line) {
		close(signal)
		//remover.Remove()
	})
}
