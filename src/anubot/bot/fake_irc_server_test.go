package bot_test

import (
	"errors"
	"net"
	"strings"
	"sync"

	. "github.com/onsi/gomega"
)

type fakeIRCServer struct {
	listener    net.Listener
	connections []net.Conn

	received  map[net.Conn][]byte
	sent      map[net.Conn][]byte
	responses map[int][][]byte

	stop    chan struct{}
	stopped chan struct{}
	connWg  sync.WaitGroup
	mu      sync.Mutex
}

func newFakeIRCServer(listener net.Listener) *fakeIRCServer {
	return &fakeIRCServer{
		listener: listener,

		received:  make(map[net.Conn][]byte),
		sent:      make(map[net.Conn][]byte),
		responses: make(map[int][][]byte),

		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (a *fakeIRCServer) Start() {
	defer close(a.stopped)
	for {
		// check for stop signal
		select {
		case <-a.stop:
			return
		default:
		}

		// block and accept a connection
		conn, _ := a.listener.Accept()

		// check for stop signal again since we've been blocked
		select {
		case <-a.stop:
			return
		default:
		}

		// track connection
		a.mu.Lock()
		a.connections = append(a.connections, conn)
		a.connWg.Add(1)
		a.mu.Unlock()

		// spin off a goroutine to handle the connection
		go a.handleConn(conn)
	}
}

func (a *fakeIRCServer) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		a.removeConnection(conn)
		a.connWg.Done()
	}()
	a.mu.Lock()
	connIndex, err := a.findConnectionIndex(conn)
	a.mu.Unlock()
	Expect(err).ToNot(HaveOccurred())
	for {
		// check for stop signal
		select {
		case <-a.stop:
			return
		default:
		}

		// create a buffer
		data := make([]byte, 256)
		// block and read into buffer
		// TODO: check err?
		i, _ := conn.Read(data)

		// check for stop signal again since we've been blocked
		select {
		case <-a.stop:
			return
		default:
		}

		// write data into received
		a.mu.Lock()
		a.received[conn] = append(a.received[conn], data[:i]...)

		// tear down connection if we got a quit command
		if strings.HasPrefix(string(data[:i]), "QUIT") {
			a.mu.Unlock()
			return
		}

		// send response
		if len(a.responses[connIndex]) > 0 {
			data, a.responses[connIndex] = a.responses[connIndex][0], a.responses[connIndex][1:]
			conn.Write(data)
			a.sent[conn] = append(a.sent[conn], data...)
		}
		a.mu.Unlock()
	}
}

func (a *fakeIRCServer) Stop() {
	// signal for listeners and connections to stop
	close(a.stop)

	// close listener
	Expect(a.listener.Close()).To(Succeed())

	// close all connections
	a.mu.Lock()
	connections := a.connections
	a.mu.Unlock()
	for _, conn := range connections {
		Expect(conn.Close()).To(Succeed())
	}

	// wait for listener to stop
	<-a.stopped

	// wait for connections to stop
	a.connWg.Wait()
}

func (a *fakeIRCServer) Respond(connIndex int, data ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	message := []byte(strings.Join(data, "\r\n") + "\r\n")
	a.responses[connIndex] = append(a.responses[connIndex], message)
}

func (a *fakeIRCServer) TriggerResponse(connIndex int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	data := make([]byte, 256)
	conn := a.connections[connIndex]
	if len(a.responses[connIndex]) > 0 {
		data, a.responses[connIndex] = a.responses[connIndex][0], a.responses[connIndex][1:]
		conn.Write(data)
		a.sent[conn] = append(a.sent[conn], data...)
	}
}

func (a *fakeIRCServer) Connections() []net.Conn {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connections
}

func (a *fakeIRCServer) Received(connIndex int) func() []byte {
	return func() []byte {
		a.mu.Lock()
		defer a.mu.Unlock()
		return a.received[a.connections[connIndex]]
	}
}

func (a *fakeIRCServer) Sent(connIndex int) func() []byte {
	return func() []byte {
		a.mu.Lock()
		defer a.mu.Unlock()
		return a.sent[a.connections[connIndex]]
	}
}

func (a *fakeIRCServer) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.received = make(map[net.Conn][]byte)
	a.sent = make(map[net.Conn][]byte)
	a.responses = make(map[int][][]byte)
}

func (a *fakeIRCServer) removeConnection(conn net.Conn) {
	a.mu.Lock()
	defer a.mu.Unlock()
	i, err := a.findConnectionIndex(conn)
	Expect(err).ToNot(HaveOccurred())
	a.connections = append(a.connections[:i], a.connections[i+1:]...)
}

func (a *fakeIRCServer) findConnectionIndex(conn net.Conn) (int, error) {
	for i, c := range a.connections {
		if c == conn {
			return i, nil
		}
	}
	return 0, errors.New("connection missing")
}
