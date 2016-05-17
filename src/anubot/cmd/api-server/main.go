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
	store := store.New(store.HomePath())
	store.InitDDL()
	bot := &bot.Bot{}
	api := api.New(store, bot)
	handler := websocket.Handler(api.Serve)

	http.Handle("/api", handler)
	fmt.Println("listening on port 12345")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
