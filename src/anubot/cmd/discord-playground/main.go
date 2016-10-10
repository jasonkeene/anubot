package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	t := os.Getenv("DISCORD_BOT_PASS")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s, err := discordgo.New("Bot " + t)
	if err != nil {
		panic(err)
	}
	// Add session calls here.
	_ = s

	<-interrupt
}
