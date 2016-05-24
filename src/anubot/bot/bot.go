package bot

import (
	"crypto/tls"
	"net"
	"strconv"

	"github.com/fluffle/goirc/client"
)

type ConnConfig struct {
	UserUsername string
	UserPassword string
	BotUsername  string
	BotPassword  string
	Host         string
	Port         int
	TLSConfig    *tls.Config
}

// Bot communicates with the IRC server and has pointers to features.
type Bot struct {
	// TODO: Implement userConn (this might not be needed)
	// userConn     *client.Conn
	botConn      *client.Conn
	connected    chan struct{}
	disconnected chan struct{}
}

// Connect establishes a connection to the IRC server.
func (b *Bot) Connect(c *ConnConfig) (chan struct{}, error) {
	// TODO: Is this idempotent?
	cfg := client.NewConfig(c.BotUsername)
	cfg.Me.Name = c.BotUsername
	cfg.Me.Ident = "anubot"
	cfg.Pass = c.BotPassword
	cfg.SSL = true
	if c.TLSConfig == nil {
		c.TLSConfig = &tls.Config{
			ServerName: c.Host,
		}
	}
	cfg.SSLConfig = c.TLSConfig
	cfg.Server = net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	b.botConn = client.Client(cfg)
	b.connected = make(chan struct{})
	b.disconnected = make(chan struct{})

	b.registerConnectEventHandler()
	b.registerDisconnectEventHandler()

	return b.disconnected, b.botConn.Connect()
}

// Disconnect tears down the connection to the IRC server.
func (b *Bot) Disconnect() {
	b.botConn.Quit()
}

// Join joins an IRC channel.
func (b *Bot) Join(channel string) {
	<-b.connected
	b.botConn.Join(channel)
}

// Send sents a raw message to the IRC server.
func (b *Bot) Send(message string) {
	b.botConn.Raw(message)
}

func (b *Bot) registerConnectEventHandler() {
	b.botConn.HandleFunc(client.CONNECTED, func(conn *client.Conn, line *client.Line) {
		close(b.connected)
	})
}

func (b *Bot) registerDisconnectEventHandler() {
	b.botConn.HandleFunc(client.DISCONNECTED, func(conn *client.Conn, line *client.Line) {
		close(b.disconnected)
	})
}
