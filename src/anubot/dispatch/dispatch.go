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
		log.Panicf("can not create dispatcher publish socket: %s", err)
	}
	for _, endpoint := range pubEndpoints {
		err = pub.Bind(endpoint)
		if err != nil {
			log.Panicf("can not bind publish socket: %s", err)
		}
	}

	push, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		log.Panicf("can not create dispatcher push socket: %s", err)
	}
	for _, endpoint := range pushEndpoints {
		err = push.Bind(endpoint)
		if err != nil {
			log.Panicf("can not bind push socket: %s", err)
		}
	}

	return &Dispatcher{
		pub:  pub,
		push: push,
	}
}

// Dispatch accepts messages for sending to consumers.
func (d *Dispatcher) Dispatch(message stream.RXMessage) {
	mb, err := json.Marshal(message)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error marshalling RXMessage: %s", err)
		return
	}

	// TODO: dedupe messages
	// TODO: figure out topic

	d.mu.Lock()
	defer d.mu.Unlock()
	_, err = d.pub.SendMessage("foo", mb)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error publishing message: %s", err)
	}
	_, err = d.push.SendBytes(mb, zmq4.DONTWAIT)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error pushing message: %s", err)
	}
}
