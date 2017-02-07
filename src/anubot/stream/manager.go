package stream

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/pebbe/zmq4"
)

// Manager manages numerous connections to stream soruces.
type Manager struct {
	pushEndpoints []string
	push          *zmq4.Socket
	dispatch      chan dispatchMessage

	mu          sync.Mutex
	twitchConns map[string]conn
	discordConn conn

	twitch TwitchUserIDFetcher
}

type dispatchMessage struct {
	topic string
	msg   RXMessage
}

type conn interface {
	send(TXMessage)
	close() error
}

// TwitchUserIDFetcher fetches the user ID for a given username.
type TwitchUserIDFetcher interface {
	UserID(username string) (userID int, err error)
}

// Option is used to configure a Mananger.
type Option func(*Manager)

// WithPushEndpoints allows you to override the default push endpoints.
func WithPushEndpoints(endpoints []string) Option {
	return func(m *Manager) {
		m.pushEndpoints = endpoints
	}
}

// NewManager creates a new manager.
func NewManager(twitch TwitchUserIDFetcher, opts ...Option) *Manager {
	m := &Manager{
		pushEndpoints: []string{"inproc://dispatch-pull"},
		dispatch:      make(chan dispatchMessage, 1000),
		twitchConns:   make(map[string]conn),
		twitch:        twitch,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.setupSockets()
	go m.run()
	return m
}

func (m *Manager) setupSockets() {
	var err error
	m.push, err = zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		log.Panicf("Manager.setupSockets: can not create manager push socket: %s", err)
	}
	for _, endpoint := range m.pushEndpoints {
		err = m.push.Connect(endpoint)
		if err != nil {
			log.Panicf("Manager.setupSockets: can not connect push socket: %s", err)
		}
	}
}

func (m *Manager) run() {
	for dispatchMsg := range m.dispatch {
		mb, err := json.Marshal(dispatchMsg.msg)
		if err != nil {
			log.Printf("Manager.run: error with marshalling message: %s", err)
			continue
		}

		_, err = m.push.SendBytes([]byte(dispatchMsg.topic), zmq4.SNDMORE)
		if err != nil {
			log.Printf("Manager.run: unable to send message frame 0: %s", err)
			continue
		}
		_, err = m.push.SendBytes(mb, 0)
		if err != nil {
			log.Printf("Manager.run: unable to send message frame 1: %s", err)
			continue
		}
	}
}

// ConnectTwitch connects to twitch and streams data to the dispatcher.
func (m *Manager) ConnectTwitch(user, pass, channel string) {
	m.mu.Lock()
	_, ok := m.twitchConns[user]
	m.mu.Unlock()
	if ok {
		return
	}

	for i := 0; i < 10; i++ {
		c, err := connectTwitch(user, pass, channel, m.dispatch, m.twitch)
		if err == nil {
			m.mu.Lock()
			defer m.mu.Unlock()
			m.twitchConns[user] = c
			return
		}
	}
	log.Print("unable to establish connection to twitch for user:", user)
}

// ConnectDiscord connects to discord and streams data to the dispatcher.
func (m *Manager) ConnectDiscord(token string) {
	m.mu.Lock()
	dc := m.discordConn
	m.mu.Unlock()
	if dc != nil {
		return
	}

	for i := 0; i < 10; i++ {
		c, err := connectDiscord(token, m.dispatch)
		if err == nil {
			m.mu.Lock()
			defer m.mu.Unlock()
			m.discordConn = c
			return
		}
	}
	log.Print("unable to establish connection to discord")
}

// DisconnectTwitch tears down a connection to twitch.
func (m *Manager) DisconnectTwitch(user string) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Print("Manager.DisconnectTwitch: disconnecting for user:", user)

	c, ok := m.twitchConns[user]
	if !ok {
		log.Print("Manager.DisconnectTwitch: user conn does not exist for twitch user:", user)
		return func() {}
	}
	err := c.close()
	delete(m.twitchConns, user)
	if err != nil {
		log.Printf("Manager.DisconnectTwitch: error occurred while disconnecting user: %s error: %s", user, err)
		return func() {}
	}
	return func() {
		// TODO: block until disconnect completed
		time.Sleep(time.Second)
	}
}

// DisconnectDiscord tears down the connection to discord.
func (m *Manager) DisconnectDiscord() func() {
	log.Print("Manager.DisconnectDiscord: disconnecting")
	m.mu.Lock()
	c := m.discordConn
	m.discordConn = nil
	m.mu.Unlock()
	err := c.close()
	if err != nil {
		log.Printf("Manager.DisconnectDiscord: error occurred while disconnecting: %s", err)
		return func() {}
	}
	return func() {
		// TODO: block until disconnect completed
		time.Sleep(time.Second)
	}
}

// Send sends a message to the stream source.
func (m *Manager) Send(ms TXMessage) {
	var c conn
	switch ms.Type {
	case Twitch:
		m.mu.Lock()
		c = m.twitchConns[ms.Twitch.Username]
		m.mu.Unlock()
		if c == nil {
			log.Printf("unable to send message for twitch user: %s", ms.Twitch.Username)
			return
		}
	case Discord:
		m.mu.Lock()
		c = m.discordConn
		m.mu.Unlock()
	default:
		log.Printf("Manager.Send: unknown message type: %d", ms.Type)
		return
	}
	c.send(ms)
}
