package main

import (
	"os"

	"anubot/store"
)

func main() {
	twitchUser := os.Getenv("TWITCH_USER")
	twitchPass := os.Getenv("TWITCH_PASS")
	twitchChannel := os.Getenv("TWITCH_CHANNEL")

	s := store.New(store.HomePath())
	s.InitDDL()
	s.SetCredentials(twitchUser, twitchPass)
	s.SetPrimaryChannel(twitchChannel)
}
