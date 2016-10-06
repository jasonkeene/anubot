package dispatch

import (
	"anubot/stream"
	"encoding/json"
	"log"
	"sync"

	"github.com/pebbe/zmq4"
)

// Dispatcher dispatches messages to consumers.
type Dispatcher struct {
	mu   sync.Mutex
	pub  *zmq4.Socket
	push *zmq4.Socket
}

// New creates a new dispatcher.
func New(pubEndpoints, pushEndpoints []string) *Dispatcher {
	pub, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		log.Panicf("dispatch.New: can not create dispatcher publish socket: %s", err)
	}
	for _, endpoint := range pubEndpoints {
		err = pub.Bind(endpoint)
		if err != nil {
			log.Panicf("dispatch.New: can not bind publish socket: %s", err)
		}
	}

	push, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		log.Panicf("dispatch.New: can not create dispatcher push socket: %s", err)
	}
	for _, endpoint := range pushEndpoints {
		err = push.Bind(endpoint)
		if err != nil {
			log.Panicf("dispatch.New: can not bind push socket: %s", err)
		}
	}

	return &Dispatcher{
		pub:  pub,
		push: push,
	}
}

// Dispatch accepts messages for sending to consumers.
func (d *Dispatcher) Dispatch(topic string, message stream.RXMessage) {
	mb, err := json.Marshal(message)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error marshalling RXMessage: %s", err)
		return
	}

	switch message.Type {
	case stream.Twitch:
		if len(message.Twitch.Line.Args) < 2 {
			log.Print("Dispatcher.Dispatch: twtich message did not have channel name")
			return
		}
	case stream.Discord:
	default:
		log.Printf("Dispatcher.Dispatch: unknown message type: %d", message.Type)
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	_, err = d.pub.SendMessage(topic, mb)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error publishing message: %s", err)
	}
	_, err = d.push.SendBytes(mb, zmq4.DONTWAIT)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error pushing message: %s", err)
	}
}
