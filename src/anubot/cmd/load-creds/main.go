package main

import (
	"os"

	"anubot/store"
)

func main() {
	botUser := os.Getenv("TWITCH_BOT_USER")
	botPass := os.Getenv("TWITCH_BOT_PASS")
	userUser := os.Getenv("TWITCH_USER_USER")
	userPass := os.Getenv("TWITCH_USER_PASS")

	s := store.New(store.HomePath())
	s.InitDDL()
	s.SetCredentials("bot", botUser, botPass)
	s.SetCredentials("user", userUser, userPass)
}
