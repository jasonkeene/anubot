package main

import (
	"bufio"
	"os"
	"os/signal"

	"github.com/davecgh/go-spew/spew"
	"github.com/fluffle/goirc/logging/golog"

	"anubot/stream"
)

func init() {
	golog.Init()
}

type fakeDispatcher struct{}

func (fd fakeDispatcher) Dispatch(m stream.RXMessage) {
	spew.Dump(m)
}

func main() {
	u := os.Getenv("TWITCH_USER_USER")
	p := os.Getenv("TWITCH_USER_PASS")

	manager := stream.NewManager()
	err := manager.Connect(stream.Twitch, u, p, "#"+u, fakeDispatcher{})
	if err != nil {
		panic(err)
	}

	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				panic(err)
			}
			manager.Send(stream.Twitch, u, stream.TXMessage{
				To:      "#" + u,
				Message: string(line),
			})
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	manager.Disconnect(stream.Twitch, u)
}
