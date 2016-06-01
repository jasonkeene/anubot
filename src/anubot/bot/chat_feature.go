package bot

import "github.com/fluffle/goirc/client"

type ChatFeature struct {
	writer     FeatureWriter
	dispatcher Dispatcher
}

func NewChatFeature(writer FeatureWriter, dispatcher Dispatcher) *ChatFeature {
	return &ChatFeature{
		writer:     writer,
		dispatcher: dispatcher,
	}
}

func (cf *ChatFeature) Init() {
	cf.writer.HandleFunc("bot", "PRIVMSG", cf.ChatHandler)
	cf.writer.HandleFunc("streamer", "PRIVMSG", cf.ChatHandler)
}

func (cf *ChatFeature) Send(user, message string) {
	cf.writer.Privmsg(user, cf.writer.Channel(), message)
}

func (cf *ChatFeature) ChatHandler(conn *client.Conn, line *client.Line) {
	// TODO: Args might not always exist
	cf.dispatcher.Dispatch(Message{
		Channel: line.Args[0],
		Body:    line.Args[1],
		Time:    line.Time,
	})
}
