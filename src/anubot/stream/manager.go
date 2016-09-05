package stream

import (
	"log"
	"sync"
)

// Manager manages numerous connections to stream soruces.
type Manager struct {
	mu    sync.Mutex
	conns map[connKey]conn
}

// NewManager creates a new manager.
func NewManager() *Manager {
	return &Manager{
		conns: make(map[connKey]conn),
	}
}

type connKey struct {
	t Type
	u string
}

type conn interface {
	send(TXMessage)
	close() error
}

// Connect establishes a connection to the stream source. If the connection
// fails initially an error is returned. If the connection fails later it will
// attempt to reconnect.
func (m *Manager) Connect(t Type, u, p, c string, d Dispatcher) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := connKey{t: t, u: u}
	switch t {
	case Twitch:
		c, err := connectTwitch(u, p, c, d)
		if err != nil {
			return err
		}
		m.conns[key] = c
		return nil
	case Discord:
		c, err := connectDiscord(u, p, c, d)
		if err != nil {
			return err
		}
		m.conns[key] = c
		return nil
	}
	log.Panicf("unknown stream type %d", t)
	return nil
}

// Disconnect tears down a connection to the stream source.
func (m *Manager) Disconnect(t Type, u string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := connKey{t: t, u: u}
	c, ok := m.conns[key]
	if !ok {
		return nil
	}
	return c.close()
}

// Send sends a message to the stream source.
func (m *Manager) Send(t Type, u string, ms TXMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := connKey{t: t, u: u}
	c, ok := m.conns[key]
	if !ok {
		return
	}
	c.send(ms)
}
