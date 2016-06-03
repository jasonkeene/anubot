package bot

import (
	"crypto/tls"
	"errors"
	"log"
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
	Flood            bool
}

// Bot communicates with the IRC server and has pointers to features.
type Bot struct {
	mu           sync.Mutex
	connected    bool
	streamerConn *client.Conn
	botConn      *client.Conn

	channelMu sync.Mutex
	channel   string

	// features
	chatFeature *ChatFeature
}

// Connect establishes two connections to the Twitch IRC server, one as the
// streamer and one as the bot. It then joins the streamer's channel.
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

	// join streamer's channel by default
	b.channelMu.Lock()
	b.channel = "#" + c.StreamerUsername
	b.channelMu.Unlock()
	b.join(b.channel)

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
func (b *Bot) Join(channel string) {
	b.channelMu.Lock()
	defer b.channelMu.Unlock()

	b.join(channel)
	b.part(channel)
	b.channel = channel
}

func (b *Bot) join(channel string) {
	b.streamerConn.Join(channel)
	b.botConn.Join(channel)
}

func (b *Bot) part(channel string) {
	b.streamerConn.Part(b.channel)
	b.botConn.Part(b.channel)
}

// Channel returns the currently active channel.
func (b *Bot) Channel() string {
	b.channelMu.Lock()
	defer b.channelMu.Unlock()
	return b.channel
}

// Send sends a raw message to the specified IRC connection.
func (b *Bot) Send(user, message string) {
	switch user {
	case "streamer":
		b.streamerConn.Raw(message)
	case "bot":
		b.botConn.Raw(message)
	default:
		log.Panicf("Bad user provided for sending message")
	}
}

// Privmsg sends a chat message to the specified IRC connection.
func (b *Bot) Privmsg(user, target, message string) {
	switch user {
	case "streamer":
		b.streamerConn.Privmsg(target, message)
	case "bot":
		b.botConn.Privmsg(target, message)
	default:
		log.Panicf("Bad user provided for sending message")
	}
}

// HandleFunc registers functions to handle IRC commands for a specific user.
func (b *Bot) HandleFunc(user, command string, handlefunc client.HandlerFunc) {
	switch user {
	case "streamer":
		b.streamerConn.HandleFunc(command, handlefunc)
	case "bot":
		b.botConn.HandleFunc(command, handlefunc)
	default:
		log.Panicf("Bad user provided for registering handlefunc")
	}
}

// TODO: cover this in tests
// InitChatFeature wires up the chat feature of the bot.
func (b *Bot) InitChatFeature(dispatcher *MessageDispatcher) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.chatFeature = NewChatFeature(b, dispatcher)
	b.chatFeature.Init()
}

// TODO: cover this in tests
// ChatFeature returns the chat feature of the bot.
func (b *Bot) ChatFeautre() *ChatFeature {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.chatFeature
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
	cfg.Flood = c.Flood
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
