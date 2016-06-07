package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"

	"anubot/api"
	"anubot/bot"
	"anubot/store"
)

func main() {
	// create and initialize database connection
	store := store.New(store.HomePath())
	store.InitDDL()

	// setup bot that communicates with the twitch IRC server
	b := &bot.Bot{}

	// create message dispatcher
	dispatcher := bot.NewMessageDispatcher()

	// wire up features
	b.InitChatFeature(dispatcher)

	// setup websocket API server
	api := api.New(store, b, dispatcher)
	http.Handle("/api", websocket.Handler(api.Serve))

	// bind websocket API
	fmt.Println("listening on port 12345")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
