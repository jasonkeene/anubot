package api

import "golang.org/x/net/websocket"

// session stores objects handlers need when responding to events.
type session struct {
	id  string
	ws  *websocket.Conn
	api *Server

	authenticated bool
	userID        string
}

// Send sends an event to the user over the websocket connection.
func (s *session) Send(e event) {
	websocket.JSON.Send(s.ws, e)
}

// Receive returns the next event from the websocket connection.
func (s *session) Receive() (event, error) {
	var e event
	err := websocket.JSON.Receive(s.ws, &e)
	return e, err
}

// Store gets the store for this session.
func (s *session) Store() Store {
	return s.api.store
}

// TwitchOauthClientID gets the oauth client ID for Twitch.
func (s *session) TwitchOauthClientID() string {
	return s.api.twitchOauthClientID
}

// SetAuthentication sets the authentication for this session.
func (s *session) SetAuthentication(id string) {
	s.authenticated = true
	s.userID = id
}

// Authenticated lets you know what user this session is authenticated as.
func (s *session) Authenticated() (string, bool) {
	return s.userID, s.authenticated
}

// Logout clears the authentication for this session.
func (s *session) Logout() {
	s.authenticated = false
	s.userID = ""
}
