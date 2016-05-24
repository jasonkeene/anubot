package api

import (
	"anubot/bot"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

//go:generate hel -t Store -o mock_store_test.go

// Store is the object the APIServer uses to persist data.
type Store interface {
	SetCredentials(kind, user, pass string) (err error)
	HasCredentials(kind string) (has bool)
	Credentials(kind string) (user string, pass string, err error)
}

//go:generate hel -t Bot -o mock_bot_test.go

// Bot is the object reposible for talking to IRC.
type Bot interface {
	Connect(connConfig *bot.ConnConfig) (err error, disconnected chan struct{})
	Disconnect()
}

// Event is the structure sent over websocket connections by both ends.
type Event struct {
	Cmd     string      `json:"cmd"`
	Payload interface{} `json:"payload"`
}

// Session stores objects handlers need when responding to events.
type Session struct {
	ws    *websocket.Conn
	store Store
	bot   Bot
}

// APIServer responds to websocket events sent from the client.
type APIServer struct {
	store Store
	bot   Bot
}

// New creates a new APIServer.
func New(store Store, bot Bot) *APIServer {
	return &APIServer{
		store: store,
		bot:   bot,
	}
}

// Serve reads off of a websocket connection and responds to events.
func (api *APIServer) Serve(ws *websocket.Conn) {
	defer ws.Close()

	session := &Session{
		ws:    ws,
		store: api.store,
		bot:   api.bot,
	}

	for {
		var event Event
		err := websocket.JSON.Receive(ws, &event)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Panic(err)
		}
		handler, ok := eventHandlers[event.Cmd]
		if !ok {
			continue
		}
		handler.HandleEvent(event, session)
	}
}
