package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/fluffle/goirc/logging/golog"
	"github.com/spf13/viper"

	"anubot/api"
	"anubot/bot"
	"anubot/dispatch"
	"anubot/store"
	"anubot/store/bolt"
	"anubot/stream"
	"anubot/twitch"
	"anubot/twitch/oauth"
)

func init() {
	golog.Init()
}

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
	backend := v.GetString("store_backend")
	var s store.Store
	switch backend {
	case "bolt":
		var err error
		s, err = bolt.New(v.GetString("store_bolt_path"))
		if err != nil {
			log.Panicf("unable to craete bolt database: %s", err)
		}
	case "dummy":
		log.Panicf("dummy store backend is not wired up")
	default:
		log.Panicf("unknown store backend: %s", backend)
	}

	// create message dispatcher
	pubEndpoints := []string{
		"inproc://pub",
	}
	pushEndpoints := []string{
		"inproc://push",
	}
	d := dispatch.New(pubEndpoints, pushEndpoints)

	// setup puller to store messages
	puller, err := store.NewPuller(s, []string{"inproc://push"})
	if err != nil {
		log.Panicf("pull not able to connect, got err: %s", err)
	}
	go puller.Start()

	// create bot manager
	bm := bot.NewManager()

	// create stream manager
	sm := stream.NewManager(d, twitch)

	// setup websocket API server
	mux := http.NewServeMux()
	api := api.New(
		bm,
		sm,
		pubEndpoints,
		s,
		twitch,
		v.GetString("twitch_oauth_client_id"),
	)
	mux.Handle("/api", websocket.Handler(api.Serve))

	// wire up oauth handler
	mux.Handle("/twitch_oauth/done", oauth.NewDoneHandler(
		v.GetString("twitch_oauth_client_id"),
		v.GetString("twitch_oauth_client_secret"),
		v.GetString("twitch_oauth_redirect_uri"),
		s,
		twitch,
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
