package api

import (
	"anubot/bot"
	"log"

	"golang.org/x/net/websocket"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

var eventHandlers map[string]EventHandler

type EventHandler interface {
	HandleEvent(event Event, session *Session)
}

func pingHandler(event Event, session *Session) {
	websocket.JSON.Send(session.ws, &Event{Cmd: "pong"})
}

func hasCredentialsSetHandler(event Event, session *Session) {
	var result bool
	kind, ok := event.Payload.(string)
	if !ok {
		return
	}
	if kind == "user" || kind == "bot" {
		result = session.store.HasCredentials(kind)
	}
	websocket.JSON.Send(session.ws, &Event{
		Cmd: "has-credentials-set",
		Payload: struct {
			Kind   string `json:"kind"`
			Result bool   `json:"result"`
		}{
			Kind:   kind,
			Result: result,
		},
	})
}

func setCredentialsHandler(event Event, session *Session) {
	data, ok := event.Payload.(map[string]interface{})
	if !ok {
		log.Println("Unable to assert type of set-credentials payload")
		return
	}

	kind, err := getStringFromMap("kind", data)
	if err != nil {
		log.Println(err)
		return
	}
	username, err := getStringFromMap("username", data)
	if err != nil {
		log.Println(err)
		return
	}
	password, err := getStringFromMap("password", data)
	if err != nil {
		log.Println(err)
		return
	}

	// validation of values
	if kind == "" || username == "" || password == "" {
		log.Println("Invalid credentials provided in set-credentials event")
		return
	}
	if kind != "user" && kind != "bot" {
		log.Println("Bad kind provided in set-credentials event")
		return
	}
	session.store.SetCredentials(kind, username, password)
}

func connectHandler(event Event, session *Session) {
	userUser, userPass, err := session.store.Credentials("user")
	if err != nil {
		return
	}
	botUser, botPass, err := session.store.Credentials("bot")
	if err != nil {
		return
	}
	connConfig := &bot.ConnConfig{
		UserUsername: userUser,
		UserPassword: userPass,
		BotUsername:  botUser,
		BotPassword:  botPass,
		Host:         twitchHost,
		Port:         twitchPort,
	}
	// TODO: handle error from connect as a false payload
	session.bot.Connect(connConfig)
	websocket.JSON.Send(session.ws, &Event{
		Cmd:     "connect",
		Payload: true,
	})
}

func disconnectHandler(event Event, session *Session) {
	session.bot.Disconnect()
}

type HandlerFunc func(event Event, session *Session)

func (f HandlerFunc) HandleEvent(event Event, session *Session) {
	log.Printf("handling %s event", event.Cmd)
	f(event, session)
}

func init() {
	eventHandlers = make(map[string]EventHandler)
	eventHandlers["ping"] = HandlerFunc(pingHandler)
	eventHandlers["has-credentials-set"] = HandlerFunc(hasCredentialsSetHandler)
	eventHandlers["set-credentials"] = HandlerFunc(setCredentialsHandler)
	eventHandlers["connect"] = HandlerFunc(connectHandler)
	eventHandlers["disconnect"] = HandlerFunc(disconnectHandler)
}
