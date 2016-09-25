package bot

import (
	"anubot/stream"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/pebbe/zmq4"
)

// Sender sends messsgaes to a stream source.
type Sender interface {
	Send(ms stream.TXMessage)
}

// Feature accepts messages and tick signals that are used to implement bot
// logic.
type Feature interface {
	HandleMessage(ms stream.RXMessage)
	Tick()
}

// Bot receives messages and takes actions based on those messages.
type Bot struct {
	sub        *zmq4.Socket
	featuresMu sync.Mutex
	features   []Feature
}

// New returns a new Bot that is connected to publishers.
func New(topic string, pubEndpoints []string) (*Bot, error) {
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

	b := &Bot{
		sub: sub,
	}
	return b, nil
}

// Start spawns goroutines needed to handle messages.
func (b *Bot) Start() {
	go b.startSub()
	go b.startTicker()
}

// Stop tears down the goroutines needed to handle messages.
func (b *Bot) Stop() {
	// TODO: implement this
}

func (b *Bot) startSub() {
	for {
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

		var fs []Feature
		b.featuresMu.Lock()
		copy(fs, b.features)
		b.featuresMu.Unlock()
		for _, f := range fs {
			f.HandleMessage(ms)
		}
	}
}

func (b *Bot) startTicker() {
	t := time.NewTicker(time.Second)
	for range t.C {
		var fs []Feature
		b.featuresMu.Lock()
		copy(fs, b.features)
		b.featuresMu.Unlock()
		for _, f := range fs {
			f.Tick()
		}
	}
}

// AddFeature adds a feature to accept messages and ticks.
func (b *Bot) AddFeature(f Feature) {
	b.featuresMu.Lock()
	b.featuresMu.Unlock()
	b.features = append(b.features, f)
}
