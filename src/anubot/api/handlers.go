package api

import (
	"anubot/bot"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}

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
		StreamerUsername: userUser,
		StreamerPassword: userPass,
		BotUsername:      botUser,
		BotPassword:      botPass,
		Host:             twitchHost,
		Port:             twitchPort,
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

func subscribeHandler(event Event, session *Session) {
	messages := session.dispatcher.Messages(session.bot.Channel())
	session.AddSubscription(messages)

	go func() {
		for msg := range messages {
			websocket.JSON.Send(session.ws, &Event{
				Cmd:     "chat-message",
				Payload: msg,
			})
		}
	}()
}

func unsubscribeHandler(event Event, session *Session) {
	// TODO: need to implement
}

func sendMessageHandler(event Event, session *Session) {
	data, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}
	user, ok := data["user"].(string)
	if !ok {
		return
	}
	message, ok := data["message"].(string)
	if !ok {
		return
	}
	session.bot.ChatFeature().Send(user, message)
}

func usernamesHandler(event Event, session *Session) {
	websocket.JSON.Send(session.ws, &Event{
		Cmd: "usernames",
		Payload: map[string]string{
			"streamer": session.bot.StreamerUsername(),
			"bot":      session.bot.BotUsername(),
		},
	})
}

func authDetailsHandler(event Event, session *Session) {
	resp := &Event{
		Cmd: "auth-details",
		Payload: map[string]interface{}{
			"authenticated": false,
			"streamer":      "",
			"bot":           "",
			"status":        "",
			"game":          "",
		},
	}
	defer websocket.JSON.Send(session.ws, resp)

	streamerUser, streamerPass, err := session.store.Credentials("user")
	if err != nil {
		fmt.Println(err)
		return
	}
	botUser, botPass, err := session.store.Credentials("bot")
	if err != nil {
		fmt.Println(err)
		return
	}
	if streamerUser == "" || streamerPass == "" ||
		botUser == "" || botPass == "" {
		println("bad creds")
		return
	}

	status, game, _ := fetchStreamInfo(streamerUser)
	resp.Payload = map[string]interface{}{
		"authenticated": true,
		"streamer":      streamerUser,
		"bot":           botUser,
		"status":        status,
		"game":          game,
	}
}

func fetchStreamInfo(channel string) (string, string, error) {
	url := fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s", channel)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	data := &struct {
		Status string `json:"status"`
		Game   string `json:"game"`
	}{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", "", err
	}
	return data.Status, data.Game, nil
}

func updateDescriptionHandler(event Event, session *Session) {
	payload := event.Payload.(map[string]interface{})
	game := payload["game"].(string)
	status := payload["status"].(string)
	user, pass, err := session.store.Credentials("user")
	if err != nil {
		fmt.Println("bad creds!")
		return
	}
	err = updateDescription(game, status, user, pass)
	if err != nil {
		fmt.Println(err)
	}
}

func updateDescription(game, status, channel, pass string) error {
	url := fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s", channel)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	// TODO: you need a oauth token that had the api scope!
	token := fmt.Sprintf("OAuth %s", pass)
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code %d", resp.StatusCode)
	}
	return nil
}

type HandlerFunc func(event Event, session *Session)

func (f HandlerFunc) HandleEvent(event Event, session *Session) {
	log.Printf("handling %s event", event.Cmd)
	f(event, session)
}

func init() {
	eventHandlers = make(map[string]EventHandler)
	eventHandlers["ping"] = HandlerFunc(pingHandler)
	// TODO: replace with auth-details call
	eventHandlers["has-credentials-set"] = HandlerFunc(hasCredentialsSetHandler)
	eventHandlers["set-credentials"] = HandlerFunc(setCredentialsHandler)
	eventHandlers["connect"] = HandlerFunc(connectHandler)
	eventHandlers["disconnect"] = HandlerFunc(disconnectHandler)
	eventHandlers["subscribe"] = HandlerFunc(subscribeHandler)
	eventHandlers["unsubscribe"] = HandlerFunc(unsubscribeHandler)
	eventHandlers["send-message"] = HandlerFunc(sendMessageHandler)
	// TODO: replace with auth-details call
	eventHandlers["usernames"] = HandlerFunc(usernamesHandler)
	eventHandlers["auth-details"] = HandlerFunc(authDetailsHandler)
	eventHandlers["update-description"] = HandlerFunc(updateDescriptionHandler)
}
