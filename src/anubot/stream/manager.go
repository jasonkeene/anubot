package stream

import (
	"log"
	"sync"
	"time"
)

// Manager manages numerous connections to stream soruces.
type Manager struct {
	d           Dispatcher
	mu          sync.Mutex
	twitchConns map[string]conn
	discordConn conn
}

type conn interface {
	send(TXMessage)
	close() error
}

// NewManager creates a new manager.
func NewManager(d Dispatcher) *Manager {
	return &Manager{
		d:           d,
		twitchConns: make(map[string]conn),
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
		c, err := connectTwitch(user, pass, channel, m.d)
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
		c, err := connectDiscord(token, m.d)
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
		log.Printf("Manager.DisconnectTwitch: error occured while disconnecting user: %s error: %s", user, err)
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
		log.Printf("Manager.DisconnectDiscord: error occured while disconnecting: %s", err)
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
