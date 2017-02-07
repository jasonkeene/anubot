package dispatch

import (
	"log"

	"github.com/pebbe/zmq4"
)

// Dispatcher receives messages and sends them to the appropriate locations.
// It is meant to be easily horizontally scalable.
//
// Messages are read from a pull socket as a two frame zmq message. The first
// frame is the topic used to publish the message. The second frame is the
// actual message data.
type Dispatcher struct {
	pullEndpoints []string
	pubEndpoints  []string
	pushEndpoints []string
	pull          *zmq4.Socket
	pub           *zmq4.Socket
	push          *zmq4.Socket
}

// Option is used to configure a Dispatcher.
type Option func(*Dispatcher)

// WithPullEndpoints allows you to override the default pull endpoints.
func WithPullEndpoints(endpoints []string) Option {
	return func(d *Dispatcher) {
		d.pullEndpoints = endpoints
	}
}

// WithPubEndpoints allows you to override the default pub endpoints.
func WithPubEndpoints(endpoints []string) Option {
	return func(d *Dispatcher) {
		d.pubEndpoints = endpoints
	}
}

// WithPushEndpoints allows you to override the default push endpoints.
func WithPushEndpoints(endpoints []string) Option {
	return func(d *Dispatcher) {
		d.pushEndpoints = endpoints
	}
}

// Start creates a new dispatcher and starts it.
func Start(opts ...Option) *Dispatcher {
	d := &Dispatcher{
		pullEndpoints: []string{"inproc://dispatch-pull"},
		pubEndpoints:  []string{"inproc://dispatch-pub"},
		pushEndpoints: []string{"inproc://dispatch-push"},
	}
	for _, opt := range opts {
		opt(d)
	}
	d.setupSockets()
	go d.run()
	return d
}

func (d *Dispatcher) setupSockets() {
	var err error
	d.pull, err = zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		log.Panicf("Dispatcher.setupSockets: can not create dispatcher pull socket: %s", err)
	}
	for _, endpoint := range d.pullEndpoints {
		err = d.pull.Bind(endpoint)
		if err != nil {
			log.Panicf("Dispatcher.setupSockets: can not bind pull socket: %s", err)
		}
	}

	d.pub, err = zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		log.Panicf("Dispatcher.setupSockets: can not create dispatcher publish socket: %s", err)
	}
	for _, endpoint := range d.pubEndpoints {
		err = d.pub.Bind(endpoint)
		if err != nil {
			log.Panicf("Dispatcher.setupSockets: can not bind publish socket: %s", err)
		}
	}

	d.push, err = zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		log.Panicf("Dispatcher.setupSockets: can not create dispatcher push socket: %s", err)
	}
	for _, endpoint := range d.pushEndpoints {
		err = d.push.Bind(endpoint)
		if err != nil {
			log.Panicf("Dispatcher.setupSockets: can not bind push socket: %s", err)
		}
	}
}

func (d *Dispatcher) run() {
	for {
		parts, err := d.pull.RecvMessageBytes(0)
		if err != nil {
			log.Printf("Dispatcher.run: error occurred when reading from pull socket: %s", err)
			continue
		}
		if len(parts) != 2 {
			log.Printf("Dispatcher.run: not the right count of parts, expected 2, was: %v", parts)
			continue
		}
		topic := parts[0]
		message := parts[1]

		_, err = d.pub.SendMessage(topic, message)
		if err != nil {
			log.Printf("Dispatcher.run: got error publishing message: %s", err)
		}
		_, err = d.push.SendBytes(message, zmq4.DONTWAIT)
		if err != nil {
			log.Printf("Dispatcher.run: got error pushing message: %s", err)
		}
	}
}
