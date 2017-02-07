package api

import (
	"io"
	"log"
	"net"

	"golang.org/x/net/websocket"

	"github.com/satori/go.uuid"

	"anubot/bot"
	"anubot/stream"
	"anubot/twitch"
	"anubot/twitch/oauth"
)

// Store is the object the Server uses to persist data.
type Store interface {
	RegisterUser(username, password string) (userID string, err error)
	AuthenticateUser(username, password string) (userID string, authenticated bool)
	TwitchClearAuth(userID string)
	TwitchAuthenticated(userID string) (authenticated bool)
	TwitchStreamerAuthenticated(userID string) (authenticated bool)
	TwitchStreamerCredentials(userID string) (string, string, int)
	TwitchBotAuthenticated(userID string) (authenticated bool)
	TwitchBotCredentials(userID string) (string, string, int)
	FetchRecentMessages(userID string) ([]stream.RXMessage, error)

	oauth.NonceStore
}

// Server responds to websocket events sent from the client.
type Server struct {
	botManager          *bot.Manager
	streamManager       *stream.Manager
	subEndpoints        []string
	store               Store
	twitchClient        *twitch.API
	twitchOauthClientID string
}

// Option is used to configure a Server.
type Option func(*Server)

// WithSubEndpoints allows you to override the default endpoints that the
// server will attempt to subscribe to.
func WithSubEndpoints(endpoints []string) Option {
	return func(s *Server) {
		s.subEndpoints = endpoints
	}
}

// New creates a new Server.
func New(
	botManager *bot.Manager,
	streamManager *stream.Manager,
	store Store,
	twitchClient *twitch.API,
	twitchOauthClientID string,
	opts ...Option,
) *Server {
	s := &Server{
		botManager:          botManager,
		streamManager:       streamManager,
		subEndpoints:        []string{"inproc://dispatch-pub"},
		store:               store,
		twitchClient:        twitchClient,
		twitchOauthClientID: twitchOauthClientID,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Serve reads off of a websocket connection and responds to events.
func (api *Server) Serve(ws *websocket.Conn) {
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Printf("got error while closing ws conn: %s", err)
		}
	}()

	s := &session{
		id:  uuid.NewV4().String(),
		ws:  ws,
		api: api,
	}
	log.Printf("Serving session %s", s.id)
	defer log.Printf("done Serving session %s", s.id)

	for {
		e, err := s.Receive()
		if err != nil {
			if err == io.EOF {
				return
			}
			if _, ok := err.(*net.OpError); ok {
				log.Print("Encountered an OpErr, tearing down connection: ", err)
				return
			}
			log.Printf("Encountered an error when trying to receive an event from a websocket connection: %T %s", err, err)
			continue
		}
		handler, ok := eventHandlers[e.Cmd]
		if !ok {
			log.Printf("Received an event with the command '%s' that does not match any of our handlers.", e.Cmd)
			s.Send(event{
				Cmd:       e.Cmd,
				RequestID: e.RequestID,
				Error:     invalidCommand,
			})
			continue
		}
		log.Printf("Handling '%s' event.", e.Cmd)
		handler.HandleEvent(e, s)
	}
}
