package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"

	"anubot/api"
	"anubot/bot"
	"anubot/dispatch"
	"anubot/store/dummy"
	"anubot/stream"
	"anubot/twitch"
	"anubot/twitch/oauth"
)

func main() {
	// load config
	v := viper.New()
	v.SetEnvPrefix("anubot")
	v.AutomaticEnv()

	// setup twitch api client
	twitch := twitch.New(
		v.GetString("twitch_api_url"),
		v.GetString("twitch_oauth_client_id"),
	)

	// create store
	store := dummy.New(twitch)

	// create message dispatcher
	pubEndpoints := []string{
		"inproc://pub",
	}
	pushEndpoints := []string{
		"inproc://push",
	}
	d := dispatch.New(pubEndpoints, pushEndpoints)

	// setup dummy pull sock to prevent the push sock from blocking
	pull, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		log.Panicf("pull not created, got err: %s", err)
	}
	err = pull.Connect("inproc://push")
	if err != nil {
		log.Panicf("pull not able to connect, got err: %s", err)
	}
	go func() {
		for {
			pull.RecvBytes(0)
		}
	}()

	// create bot manager
	bm := bot.NewManager()

	// create stream manager
	sm := stream.NewManager(d)

	// setup websocket API server
	mux := http.NewServeMux()
	api := api.New(
		bm,
		sm,
		[]string{},
		store,
		twitch,
		v.GetString("twitch_oauth_client_id"),
	)
	mux.Handle("/api", websocket.Handler(api.Serve))

	// wire up oauth handler
	mux.Handle("/twitch_oauth/done", oauth.NewDoneHandler(
		v.GetString("twitch_oauth_client_id"),
		v.GetString("twitch_oauth_client_secret"),
		v.GetString("twitch_oauth_redirect_uri"),
		store,
	))

	// bind websocket API
	v.SetDefault("port", 443)
	port := v.GetInt("port")
	fmt.Println("listening on port", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		// TODO: consider timeouts
	}
	err = server.ListenAndServeTLS(
		v.GetString("tls_cert_file"),
		v.GetString("tls_key_file"),
	)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
