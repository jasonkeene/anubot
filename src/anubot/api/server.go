package api

import (
	"io"
	"log"

	"golang.org/x/net/websocket"
)

var eventHandlers map[string]EventHandler

type EventHandler interface {
	HandleEvent(event Event, ws *websocket.Conn, store Store)
}

type Event struct {
	Cmd     string      `json:"cmd"`
	Payload interface{} `json:"payload"`
}

//go:generate hel -t Store -o mock_store_test.go

type Store interface {
	HasCredentials() (has bool)
	SetCredentials(user, pass string) (err error)
}

type APIServer struct {
	store Store
}

func New(store Store) *APIServer {
	return &APIServer{
		store: store,
	}
}

func (api *APIServer) Serve(ws *websocket.Conn) {
	defer ws.Close()

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
		handler.HandleEvent(event, ws, api.store)
	}
}

type HandlerFunc func(event Event, ws *websocket.Conn, store Store)

func (f HandlerFunc) HandleEvent(event Event, ws *websocket.Conn, store Store) {
	f(event, ws, store)
}

func pingHandler(event Event, ws *websocket.Conn, store Store) {
	websocket.JSON.Send(ws, &Event{Cmd: "pong"})
}

func hasCredentialsSetHandler(event Event, ws *websocket.Conn, store Store) {
	websocket.JSON.Send(ws, &Event{
		Cmd:     "has-credentials-set",
		Payload: store.HasCredentials(),
	})
}

func setCredentialsHandler(event Event, ws *websocket.Conn, store Store) {
	data, ok := event.Payload.(map[string]interface{})
	if !ok {
		log.Println("Unable to assert type of set-credentials payload")
		return
	}

	userData, ok := data["username"]
	if !ok {
		log.Println("Username not provided in set-credentials event")
		return
	}
	user, ok := userData.(string)
	if !ok {
		log.Println("Unable to assert type of username in set-credentials event")
		return
	}
	if user == "" {
		log.Println("Empty username provided in set-credentials event")
		return
	}

	passData, ok := data["password"]
	if !ok {
		log.Println("Password not provided in set-credentials event")
		return
	}
	pass, ok := passData.(string)
	if !ok {
		log.Println("Unable to assert type of password in set-credentials event")
		return
	}
	if pass == "" {
		log.Println("Empty password provided in set-credentials event")
		return
	}

	store.SetCredentials(user, pass)
}

func init() {
	eventHandlers = make(map[string]EventHandler)
	eventHandlers["ping"] = HandlerFunc(pingHandler)
	eventHandlers["has-credentials-set"] = HandlerFunc(hasCredentialsSetHandler)
	eventHandlers["set-credentials"] = HandlerFunc(setCredentialsHandler)
}
