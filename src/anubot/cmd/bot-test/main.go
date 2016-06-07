package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		StreamerUsername: os.Getenv("TWITCH_USER_USER"),
		StreamerPassword: os.Getenv("TWITCH_USER_PASS"),
		BotUsername:      os.Getenv("TWITCH_BOT_USER"),
		BotPassword:      os.Getenv("TWITCH_BOT_PASS"),
		Host:             twitchHost,
		Port:             twitchPort,
	}

	// setup bot that communicates with the twitch IRC server
	b := &bot.Bot{}

	// create message dispatcher and subscribe to channel
	dispatcher := bot.NewMessageDispatcher()
	chanMessages := dispatcher.Messages("#" + connConfig.StreamerUsername)

	// wire up features
	b.InitChatFeature(dispatcher)

	// connect bot
	disconnected, err := b.Connect(connConfig)
	if err != nil {
		panic(err)
	}

	// read from stdin and send to irc server
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			text = strings.Trim(text, "\r\n")
			if len(text) > 0 {
				b.Send("bot", text)
			}
		}
	}()
	go func() {
		for {
			msg := <-chanMessages
			fmt.Println("got message:", msg.Body)
		}
	}()
	<-disconnected
}

func initLogging() {
	logger := logginglib.NewFromFlags()
	logger.SetLogLevel(logginglib.LogDebug)
	logging.SetLogger(logger)
}
