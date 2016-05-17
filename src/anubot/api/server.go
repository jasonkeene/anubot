package api

import (
	"anubot/bot"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

//go:generate hel -t Store -o mock_store_test.go

type Store interface {
	SetCredentials(kind, user, pass string) (err error)
	HasCredentials(kind string) (has bool)
	Credentials(kind string) (user string, pass string, err error)
}

//go:generate hel -t Bot -o mock_bot_test.go

type Bot interface {
	Connect(connConfig *bot.ConnConfig) (err error, disconnected chan struct{})
	Disconnect()
}

type Event struct {
	Cmd     string      `json:"cmd"`
	Payload interface{} `json:"payload"`
}

type Session struct {
	ws    *websocket.Conn
	store Store
	bot   Bot
}

type APIServer struct {
	store Store
	bot   Bot
}

func New(store Store, bot Bot) *APIServer {
	return &APIServer{
		store: store,
		bot:   bot,
	}
}

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
