package bot

import "sync"

// Manager keeps track of active bots.
type Manager struct {
	mu sync.Mutex
	m  map[string]*Bot
}

// NewManager cerates a new manager.
func NewManager() *Manager {
	return &Manager{
		m: make(map[string]*Bot),
	}
}

// SetBot sets a bot for a given user ID.
func (m *Manager) SetBot(userID string, b *Bot) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[userID] = b
}

// GetBot gets a bot for a given user ID.
func (m *Manager) GetBot(userID string) *Bot {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m[userID]
}

// RemoveBot removes a bot for a given user ID.
func (m *Manager) RemoveBot(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, userID)
}

// Absent will only run the provided func if the bot for userID is absent.
func (m *Manager) Absent(userID string, f func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.m[userID]
	if !ok {
		f()
	}
}
