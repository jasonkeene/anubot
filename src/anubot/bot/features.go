package bot

import "github.com/fluffle/goirc/client"

//go:generate hel --type Dispatcher --output mock_dispatcher_test.go

type Dispatcher interface {
	Dispatch(msg Message)
}

//go:generate hel --type FeatureWriter --output mock_feature_writer_test.go

type FeatureWriter interface {
	HandleFunc(user, command string, handlefunc client.HandlerFunc)
	Privmsg(user, target, msg string)
	Channel() (channel string)
	StreamerUsername() string
	BotUsername() string
}
