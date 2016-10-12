package api

import (
	"anubot/stream"
	"encoding/json"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/pebbe/zmq4"
	"golang.org/x/net/websocket"
)

type message struct {
	Type    stream.Type     `json:"type"`
	Twitch  *twitchMessage  `json:"twitch"`
	Discord *discordMessage `json:"discord"`
}

type twitchMessage struct {
	Nick string            `json:"nick"`
	Body string            `json:"body"`
	Time time.Time         `json:"time"`
	Tags map[string]string `json:"tags"`
}

type discordMessage struct{}

type messageWriter struct {
	sub *zmq4.Socket
	ws  *websocket.Conn
}

func newMessageWriter(
	topic string,
	pubEndpoints []string,
	ws *websocket.Conn,
) (*messageWriter, error) {
	println("creating message writer for topic:", topic)
	fmt.Printf("pub endpoints: %v", pubEndpoints)
	sub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return nil, err
	}
	for _, endpoint := range pubEndpoints {
		err = sub.Connect(endpoint)
		if err != nil {
			return nil, err
		}
	}
	err = sub.SetSubscribe(topic)
	if err != nil {
		return nil, err
	}
	return &messageWriter{
		sub: sub,
		ws:  ws,
	}, nil
}

// start reads messages off its sub sock and writes them to its ws conn
func (mw *messageWriter) start() {
	println("starting message writer")
	for {
		rb, err := mw.sub.RecvMessageBytes(0) //zmq4.DONTWAIT)
		println("something read")
		if err != nil {
			if zmq4.AsErrno(err) != zmq4.Errno(syscall.EAGAIN) {
				log.Printf("messages not read, got err: %s", err)
			}
			continue
		}
		if len(rb) < 2 {
			log.Printf("received message bytes had invalid length: %#v", rb)
			continue
		}

		println("got message")

		var ms stream.RXMessage
		err = json.Unmarshal(rb[1], &ms)
		if err != nil {
			log.Printf("could not unmarshal, got err: %s", err)
			continue
		}

		p := message{
			Type: ms.Type,
		}
		switch ms.Type {
		case stream.Twitch:
			p.Twitch = &twitchMessage{
				Nick: ms.Twitch.Line.Nick,
				Body: ms.Twitch.Line.Args[1],
				Time: ms.Twitch.Line.Time,
				Tags: ms.Twitch.Line.Tags,
			}
		case stream.Discord:
			// TODO: add support for discord messages
			continue
		default:
			log.Println("got unknow message type while reading from sub sock")
			continue
		}
		e := event{
			Cmd:     "chat-message",
			Payload: p,
		}
		err = websocket.JSON.Send(mw.ws, e)
		if err != nil {
			log.Printf("got error when writing to ws conn, aborting: %s", err)
			return
		}
	}
}
