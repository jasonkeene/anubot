package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/spf13/viper"

	"anubot/api"
	"anubot/store/dummy"
	"anubot/twitch"
	"anubot/twitch/oauth"
)

func main() {
	// load config
	v := viper.New()
	v.SetEnvPrefix("anubot")
	v.AutomaticEnv()

	// setup twitch api for HTTP calls
	twitch := twitch.New(v.GetString("twitch_api_url"))

	// create and initialize database connection
	s := dummy.New(twitch)

	// setup websocket API server
	mux := http.NewServeMux()
	api := api.New(v.GetString("twitch_oauth_client_id"), s)
	mux.Handle("/api", websocket.Handler(api.Serve))

	// wire up oauth handler
	mux.Handle("/twitch_oauth/done", oauth.NewDoneHandler(
		v.GetString("twitch_oauth_client_id"),
		v.GetString("twitch_oauth_client_secret"),
		v.GetString("twitch_oauth_redirect_uri"),
		s,
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
	err := server.ListenAndServeTLS(
		v.GetString("tls_cert_file"),
		v.GetString("tls_key_file"),
	)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
