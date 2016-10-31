package oauth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anubot/store"
	"anubot/twitch"
)

// NonceStore is used to to store and operate on oauth nonces.
type NonceStore interface {
	CreateOauthNonce(userID string, tu store.TwitchUser) (nonce string, err error)
	OauthNonceExists(nonce string) (exists bool)
	FinishOauthNonce(nonce, username string, userID int, od store.OauthData) (err error)
}

const (
	twitchBaseURL = "https://api.twitch.tv/kraken/"
	authorizeURL  = twitchBaseURL + "oauth2/authorize"
	tokenURL      = twitchBaseURL + "oauth2/token"
	redirectURL   = "https://anubot.io/twitch_oauth/done"
	scopes        = "" +
		"user_read " +
		"user_blocks_edit " +
		"user_blocks_read " +
		"user_follows_edit " +
		"channel_read " +
		"channel_editor " +
		"channel_commercial " +
		"channel_stream " +
		"channel_subscriptions " +
		"user_subscriptions " +
		"channel_check_subscription " +
		"chat_login " +
		"channel_feed_read " +
		"channel_feed_edit"
)

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}

func parseOauthData(data []byte) (store.OauthData, error) {
	var od store.OauthData
	err := json.Unmarshal(data, &od)
	return od, err
}

// DoneHandler is where the redirect URL hits to finsih the Oauth flow.
type DoneHandler struct {
	twitchOauthClientID     string
	twitchOauthClientSecret string
	twitchOauthRedirectURI  string
	ns                      NonceStore
	twitch                  *twitch.API
}

// NewDoneHandler creates a new handler to finish the Oauth flow.
func NewDoneHandler(
	twitchOauthClientID,
	twitchOauthClientSecret,
	twitchOauthRedirectURI string,
	ns NonceStore,
	twitch *twitch.API,
) DoneHandler {
	return DoneHandler{
		twitchOauthClientID:     twitchOauthClientID,
		twitchOauthClientSecret: twitchOauthClientSecret,
		twitchOauthRedirectURI:  twitchOauthRedirectURI,
		ns:     ns,
		twitch: twitch,
	}
}

// ServeHTTP handles the response from Twitch after authentication has
// happened.
func (h DoneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	// validate nonce
	nonce := values.Get("state")
	if !h.ns.OauthNonceExists(nonce) {
		log.Print("bad nonce")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate code
	code := values.Get("code")
	if code == "" {
		log.Print("code not set in oauth response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create request to send to twitch
	r, err := http.NewRequest("POST", tokenURL, h.createPayload(nonce, code))
	if err != nil {
		log.Printf("unable to create request for posting to twitch oauth for token: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// make request to twitch
	resp, err := httpClient.Do(r)
	if err != nil {
		log.Printf("error in response from post to twitch oauth for token: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate response code
	if resp.StatusCode != http.StatusOK {
		log.Printf("got %d response code from post to twitch oauth for token", resp.StatusCode)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read response body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("DoneHandler got err in closing resp body: %s", err)
		}
	}()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("unable to read response body from post to twitch oauth for token: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse out the oauth data
	od, err := parseOauthData(d)
	if err != nil {
		log.Printf("unable to parse response from post to twitch oauth for token: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// send oauth data to the store
	user, err := h.twitch.User(od.AccessToken)
	if err != nil {
		log.Printf("unable to get user data from access token: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ns.FinishOauthNonce(nonce, user.Name, user.ID, od)
	if err != nil {
		log.Printf("unable finish oauth: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h DoneHandler) createPayload(nonce, code string) io.Reader {
	payload := url.Values{}
	payload.Set("client_id", h.twitchOauthClientID)
	payload.Set("client_secret", h.twitchOauthClientSecret)
	payload.Set("redirect_uri", h.twitchOauthRedirectURI)
	payload.Set("grant_type", "authorization_code")
	payload.Set("code", code)
	payload.Set("state", nonce)
	return strings.NewReader(payload.Encode())
}

// GenerateNonce generates a random nonce to be used in the oauth flow.
func GenerateNonce() string {
	var err error
	b := make([]byte, 20)
	for i := 0; i < 5; i++ {
		_, err = rand.Read(b)
		if err == nil {
			break
		}
	}
	if err != nil {
		panic("not able to generate a 20 byte random nonce for oauth")
	}
	return fmt.Sprintf("%x", b)
}

// URL returns a URL that will start the oauth flow.
func URL(clientID, userID string, tu store.TwitchUser, ns NonceStore) (string, error) {
	nonce, err := ns.CreateOauthNonce(userID, tu)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("redirect_uri", redirectURL)
	v.Set("scope", scopes)
	v.Set("client_id", clientID)
	v.Set("state", nonce)

	return authorizeURL + "?" + v.Encode(), nil
}
