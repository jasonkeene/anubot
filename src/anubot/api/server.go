package api

import (
	"anubot/bot"
	"io"
	"log"
	"sync"

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
	Connect(connConfig *bot.ConnConfig) (disconnected chan struct{}, err error)
	Disconnect()
	Channel() string
	InitChatFeature(dispatcher *bot.MessageDispatcher)
}

// Event is the structure sent over websocket connections by both ends.
type Event struct {
	Cmd     string      `json:"cmd"`
	Payload interface{} `json:"payload"`
}

// Session stores objects handlers need when responding to events.
type Session struct {
	ws         *websocket.Conn
	store      Store
	bot        Bot
	dispatcher *bot.MessageDispatcher

	mu            sync.Mutex
	subscriptions []chan bot.Message
}

// AddSubscription adds a subscription to the session.
func (s *Session) AddSubscription(sub chan bot.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscriptions = append(s.subscriptions, sub)
}

// Subscriptions gets the current subscriptions for the session.
func (s *Session) Subscriptions() []chan bot.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subscriptions
}

// APIServer responds to websocket events sent from the client.
type APIServer struct {
	store      Store
	bot        Bot
	dispatcher *bot.MessageDispatcher
}

// New creates a new APIServer.
func New(store Store, bot Bot, dispatcher *bot.MessageDispatcher) *APIServer {
	return &APIServer{
		store:      store,
		bot:        bot,
		dispatcher: dispatcher,
	}
}

// Serve reads off of a websocket connection and responds to events.
func (api *APIServer) Serve(ws *websocket.Conn) {
	var session *Session

	defer func() {
		// clean up session subscriptions
		for _, sub := range session.Subscriptions() {
			api.dispatcher.Remove(sub)
		}
		// disconnect from irc
		api.bot.Disconnect()
		// tear down the conn
		ws.Close()
	}()

	session = &Session{
		ws:         ws,
		store:      api.store,
		bot:        api.bot,
		dispatcher: api.dispatcher,
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
