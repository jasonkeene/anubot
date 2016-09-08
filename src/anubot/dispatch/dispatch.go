package dispatch

import (
	"anubot/stream"
	"encoding/json"
	"log"
	"sync"

	"github.com/pebbe/zmq4"
)

// Store is used to resolve users by channel name or ID,
type Store interface {
	TwitchUser(channelName string) (userID string, err error)
	DiscordUsers(channelID string) (userID []string)
}

// Dispatcher dispatches messages to consumers.
type Dispatcher struct {
	store Store

	mu   sync.Mutex
	pub  *zmq4.Socket
	push *zmq4.Socket

	recentMsgIDsMu sync.Mutex
	recentMsgIDs   *recent
}

// New creates a new dispatcher.
func New(pubEndpoints, pushEndpoints []string, store Store) *Dispatcher {
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
		store:        store,
		pub:          pub,
		push:         push,
		recentMsgIDs: newRecent(10000),
	}
}

// Dispatch accepts messages for sending to consumers.
func (d *Dispatcher) Dispatch(message stream.RXMessage) {
	mb, err := json.Marshal(message)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error marshalling RXMessage: %s", err)
		return
	}

	var topics []string
	switch message.Type {
	case stream.Twitch:
		id, ok := message.Twitch.Line.Tags["id"]
		if !ok {
			log.Print("Dispatcher.Dispatch: twitch message tags did not have id")
			return
		}
		if len(message.Twitch.Line.Args) < 2 {
			log.Print("Dispatcher.Dispatch: twtich message did not have channel name")
			return
		}
		channelName := message.Twitch.Line.Args[0]

		d.recentMsgIDsMu.Lock()
		if d.recentMsgIDs.lookup(id) {
			return
		}
		d.recentMsgIDs.insert(id)
		d.recentMsgIDsMu.Unlock()
		user, err := d.store.TwitchUser(channelName)
		if err != nil {
			log.Printf("Dispatcher.Dispatch: got error trying to resolve twitch user: %s", err)
		}
		topics = []string{user}
	case stream.Discord:
		topics = d.store.DiscordUsers(message.Discord.MessageCreate.ChannelID)
	default:
		log.Printf("Dispatcher.Dispatch: unknown message type: %d", message.Type)
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, topic := range topics {
		_, err = d.pub.SendMessage(topic, mb)
		if err != nil {
			log.Printf("Dispatcher.Dispatch: got error publishing message: %s", err)
		}
	}
	_, err = d.push.SendBytes(mb, zmq4.DONTWAIT)
	if err != nil {
		log.Printf("Dispatcher.Dispatch: got error pushing message: %s", err)
	}
}
