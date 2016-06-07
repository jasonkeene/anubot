package bot

import (
	"sync"
	"time"
)

const minBufferSize = 20

type Message struct {
	Target string    `json:"target"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
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
	// TODO: deduplicate messages that have the same channel/body
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
