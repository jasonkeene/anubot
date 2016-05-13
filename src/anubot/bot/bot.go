package anubot

import (
	"crypto/tls"
	"net"
	"strconv"

	"github.com/fluffle/goirc/client"
)

type Bot struct {
	username string
	password string

	host string
	port int

	conn         *client.Conn
	connected    chan struct{}
	disconnected chan struct{}
}

func New(username, password, host string, port int) *Bot {
	return &Bot{
		username: username,
		password: password,
		host:     host,
		port:     port,
	}
}

func (b *Bot) Connect(tlsConfig *tls.Config) (error, chan struct{}) {
	cfg := client.NewConfig(b.username)
	cfg.Me.Name = b.username
	cfg.Me.Ident = "anubot"
	cfg.Pass = b.password
	cfg.SSL = true
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName: b.host,
		}
	}
	cfg.SSLConfig = tlsConfig
	cfg.Server = net.JoinHostPort(b.host, strconv.Itoa(b.port))
	b.conn = client.Client(cfg)
	b.connected = make(chan struct{})
	b.disconnected = make(chan struct{})

	b.registerConnectEventHandler()
	b.registerDisconnectEventHandler()

	return b.conn.Connect(), b.disconnected
}

func (b *Bot) Disconnect() {
	b.conn.Quit()
}

func (b *Bot) Join(channel string) {
	<-b.connected
	b.conn.Join(channel)
}

func (b *Bot) registerConnectEventHandler() {
	b.conn.HandleFunc(client.CONNECTED, func(conn *client.Conn, line *client.Line) {
		close(b.connected)
	})
}

func (b *Bot) registerDisconnectEventHandler() {
	b.conn.HandleFunc(client.DISCONNECTED, func(conn *client.Conn, line *client.Line) {
		close(b.disconnected)
	})
}
