package main

import (
	"os"

	"github.com/fluffle/goirc/logging"
	logginglib "github.com/fluffle/golog/logging"

	"anubot/bot"
)

const (
	twitchHost = "irc.chat.twitch.tv"
	twitchPort = 443
)

func main() {
	initLogging()
	connConfig := &bot.ConnConfig{
		UserUsername: os.Getenv("TWITCH_USER_USER"),
		UserPassword: os.Getenv("TWITCH_USER_PASS"),
		BotUsername:  os.Getenv("TWITCH_BOT_USER"),
		BotPassword:  os.Getenv("TWITCH_BOT_PASS"),
		Host:         twitchHost,
		Port:         twitchPort,
	}
	bot := &bot.Bot{}
	disconnected, err := bot.Connect(connConfig)
	if err != nil {
		panic(err)
	}
	bot.Join("#" + connConfig.UserUsername)

	<-disconnected
}

func initLogging() {
	logger := logginglib.NewFromFlags()
	logger.SetLogLevel(logginglib.LogDebug)
	logging.SetLogger(logger)
}
