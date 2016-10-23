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

//go:generate hel

// Store is the object the Server uses to persist data.
type Store interface {
	RegisterUser(username, password string) (userID string, err error)
	AuthenticateUser(username, password string) (userID string, authenticated bool)
	TwitchClearAuth(userID string)
	TwitchAuthenticated(userID string) (authenticated bool)
	TwitchStreamerAuthenticated(userID string) (authenticated bool)
	TwitchStreamerCredentials(userID string) (string, string)
	TwitchBotAuthenticated(userID string) (authenticated bool)
	TwitchBotCredentials(userID string) (string, string)

	oauth.NonceStore
}

// Server responds to websocket events sent from the client.
type Server struct {
	bm                  *bot.Manager
	sm                  *stream.Manager
	pubEndpoints        []string
	store               Store
	twitch              twitch.API
	twitchOauthClientID string
}

// New creates a new Server.
func New(
	bm *bot.Manager,
	sm *stream.Manager,
	pubEndpoints []string,
	store Store,
	twitch twitch.API,
	twitchOauthClientID string,
) *Server {
	return &Server{
		bm:                  bm,
		sm:                  sm,
		pubEndpoints:        pubEndpoints,
		store:               store,
		twitch:              twitch,
		twitchOauthClientID: twitchOauthClientID,
	}
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
				Cmd:   e.Cmd,
				Error: invalidCommand,
			})
			continue
		}
		log.Printf("Handling '%s' event.", e.Cmd)
		handler.HandleEvent(e, s)
	}
}
