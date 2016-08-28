package api

import (
	"io"
	"log"
	"net"

	"github.com/satori/go.uuid"

	"golang.org/x/net/websocket"
)

//go:generate hel

// Store is the object the Server uses to persist data.
type Store interface {
	RegisterUser(username, password string) (string, error)
	AuthenticateUser(username, password string) (string, bool)

	CreateOauthNonce() string
	OauthNonceExists(nonce string) bool
}

// Server responds to websocket events sent from the client.
type Server struct {
	twitchOauthClientID string
	store               Store
}

// New creates a new Server.
func New(twitchOauthClientID string, store Store) *Server {
	return &Server{
		twitchOauthClientID: twitchOauthClientID,
		store:               store,
	}
}

// Serve reads off of a websocket connection and responds to events.
func (api *Server) Serve(ws *websocket.Conn) {
	defer func() {
		ws.Close()
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
