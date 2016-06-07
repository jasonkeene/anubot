package bot

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

const minBufferSize = 20

type Message struct {
	Nick   string    `json:"nick"`
	Target string    `json:"target"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
	ID     string    `json:"id"`
}

type MessageDispatcher struct {
	storeMu sync.Mutex
	store   map[string][]Message
	readers map[string][]chan Message
}

func NewMessageDispatcher() *MessageDispatcher {
	return &MessageDispatcher{
		store:   make(map[string][]Message),
		readers: make(map[string][]chan Message),
	}
}

func (d *MessageDispatcher) Dispatch(msg Message) {
	d.storeMu.Lock()
	defer d.storeMu.Unlock()

	d.store[msg.Target] = append(d.store[msg.Target], msg)
	for _, ch := range d.readers[msg.Target] {
		ch <- msg
	}
}

func (d *MessageDispatcher) Messages(channel string) chan Message {
	d.storeMu.Lock()
	defer d.storeMu.Unlock()

	store := d.store[channel]
	length := len(store)
	msgs := make(chan Message, max(length, minBufferSize))
	d.readers[channel] = append(d.readers[channel], msgs)
	for i := 0; i < length; i++ {
		msgs <- store[i]
	}
	return msgs
}

func (d *MessageDispatcher) Remove(messages chan Message) {
	var (
		channel string
		index   int
	)
	// TODO: this is bad code, make this not O(N^2)
loop:
	for k, readers := range d.readers {
		for i, ch := range readers {
			if ch == messages {
				index = i
				channel = k
				break loop
			}
		}
	}
	d.readers[channel] = append(d.readers[channel][:index], d.readers[channel][index+1:]...)
	close(messages)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// TODO: cover this
// WriteMessageID writes the ID field for a given message.
func WriteMessageID(message *Message) {
	message.ID = ""
	// TODO: handle error case
	mBytes, _ := json.Marshal(message)
	h := sha1.New()
	// TODO: handle error case
	h.Write(mBytes)
	message.ID = hex.EncodeToString(h.Sum(nil))
}
