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
	// Create and initialize database connection.
	store := store.New(store.HomePath())
	store.InitDDL()

	// Setup bot that communicates with the Twitch IRC server.
	b := &bot.Bot{}

	// Create message dispatcher.
	dispatcher := bot.NewMessageDispatcher()

	// Setup websocket API server.
	api := api.New(store, b, dispatcher)
	http.Handle("/api", websocket.Handler(api.Serve))

	// Bind websocket API.
	fmt.Println("listening on port 12345")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
