package api

import (
	"io"
	"log"
	"net"

	"golang.org/x/net/websocket"
)

//go:generate hel

// Store is the object the APIServer uses to persist data.
type Store interface {
	RegisterUser(username, password string) (string, error)
	AuthenticateUser(username, password string) (string, bool)

	CreateOauthNonce() string
	OauthNonceExists(nonce string) bool
}

// APIServer responds to websocket events sent from the client.
type APIServer struct {
	twitchOauthClientID string
	store               Store
}

// New creates a new APIServer.
func New(twitchOauthClientID string, store Store) *APIServer {
	return &APIServer{
		twitchOauthClientID: twitchOauthClientID,
		store:               store,
	}
}

// Serve reads off of a websocket connection and responds to events.
func (api *APIServer) Serve(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()

	s := &session{
		ws:  ws,
		api: api,
	}

	for {
		event, err := s.Receive()
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
		handler, ok := eventHandlers[event.Cmd]
		if !ok {
			log.Printf("Received an event with the command '%s' that does not match any of our handlers.", event.Cmd)
			s.Send(event{
				Cmd:   event.Cmd,
				Error: invalidCommand,
			})
			continue
		}
		log.Printf("Handling '%s' event.", event.Cmd)
		handler.HandleEvent(event, s)
	}
}
