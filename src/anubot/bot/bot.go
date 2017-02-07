package bot

import (
	"anubot/stream"
	"encoding/json"
	"log"
	"sync"

	"github.com/pebbe/zmq4"
)

// Sender sends messsgaes to a stream source.
type Sender interface {
	Send(ms stream.TXMessage)
}

// Feature accepts messages and spawns goroutines to implement the logic of
// the bot.
type Feature interface {
	HandleMessage(ms stream.RXMessage)
	Start()
	Stop()
}

// Bot receives messages and takes actions based on those messages.
type Bot struct {
	subEndpoints []string
	topics       []string
	sub          *zmq4.Socket
	featuresMu   sync.Mutex
	features     map[string]Feature
	stop         chan struct{}
	done         chan struct{}
}

// Option is used to configure a Bot.
type Option func(*Bot)

// WithSubEndpoints allows you to override the default endpoints that the
// server will attempt to subscribe to.
func WithSubEndpoints(endpoints []string) Option {
	return func(b *Bot) {
		b.subEndpoints = endpoints
	}
}

// New returns a new Bot that is connected to publishers and accepting messages
// for specific topics.
func New(topics []string, opts ...Option) (*Bot, error) {
	b := &Bot{
		subEndpoints: []string{"inproc://dispatch-pub"},
		topics:       topics,
		features:     make(map[string]Feature),
		stop:         make(chan struct{}),
		done:         make(chan struct{}),
	}
	for _, opt := range opts {
		opt(b)
	}
	err := b.setupSockets()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Bot) setupSockets() error {
	sub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return err
	}
	for _, endpoint := range b.subEndpoints {
		err = sub.Connect(endpoint)
		if err != nil {
			return err
		}
	}
	for _, topic := range b.topics {
		err = sub.SetSubscribe(topic)
		if err != nil {
			return err
		}
	}
	b.sub = sub
	return nil
}

// Start reads from sub socket and sends messages to features. It needs to run
// in its own goroutine.
func (b *Bot) Start() {
	defer close(b.done)

	for {
		select {
		case <-b.stop:
			return
		default:
		}

		rb, err := b.sub.RecvMessageBytes(0)
		if err != nil {
			log.Printf("messages not read, got err: %s", err)
			continue
		}
		if len(rb) < 2 {
			log.Printf("received message bytes had invalid length: %#v", rb)
			continue
		}
		var ms stream.RXMessage
		err = json.Unmarshal(rb[1], &ms)
		if err != nil {
			log.Printf("could not unmarshal, got err: %s", err)
			continue
		}

		b.featuresMu.Lock()
		fs := make([]Feature, 0, len(b.features))
		for _, f := range b.features {
			fs = append(fs, f)
		}
		b.featuresMu.Unlock()
		for _, f := range fs {
			f.HandleMessage(ms)
		}
	}
}

// Stop signals to the goroutine reading messages to stop. It returns a
// function that can be used to block until reading has finished.
func (b *Bot) Stop() (wait func()) {
	close(b.stop)
	return func() {
		<-b.done
	}
}

// SetFeature sets a feature to accept messages and ticks. This will overwrite
// features previously set with the same name.
func (b *Bot) SetFeature(name string, f Feature) {
	b.featuresMu.Lock()
	defer b.featuresMu.Unlock()
	b.features[name] = f
}

// RemoveFeature removes a feature from the bot and returns it.
func (b *Bot) RemoveFeature(name string) Feature {
	b.featuresMu.Lock()
	defer b.featuresMu.Unlock()
	f := b.features[name]
	delete(b.features, name)
	return f
}
