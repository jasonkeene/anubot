package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"
)

const twitchAPIURL = "https://api.twitch.tv/kraken"

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}

// API makes requests to Twitch's API.
type API struct {
	url      string
	clientID string
	mu       sync.Mutex
	games    []Game
}

// New creates a new API.
func New(url, clientID string) *API {
	if url == "" {
		url = twitchAPIURL
	}
	return &API{
		url:      url,
		clientID: clientID,
	}
}

// Username gets the username for a give oauth token.
func (t *API) Username(token string) (username string, err error) {
	u := t.url + "/user"

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Set("Client-ID", t.clientID)
	req.Header.Set("Authorization", "OAuth "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Bad status code %d", resp.StatusCode)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("got error in closing response body: %s", err)
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var userData struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(data, &userData)
	if err != nil {
		return "", err
	}
	if userData.Name == "" {
		return "", errors.New("Empty username response from twitch")
	}
	return userData.Name, nil
}

// StreamInfo returns the status and game for a given channel.
func (t *API) StreamInfo(channel string) (string, string, error) {
	u := t.url + "/channels/" + channel

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Set("Client-ID", t.clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Bad status code %d", resp.StatusCode)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("got error in closing response body: %s", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	data := &struct {
		Status string `json:"status"`
		Game   string `json:"game"`
	}{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", "", err
	}
	return data.Status, data.Game, nil
}

// UpdateDescription updates the status and game for the given channel.
func (t *API) UpdateDescription(status, game, channel, token string) error {
	u := t.url + "/channels/" + channel

	data, err := json.Marshal(map[string]map[string]string{
		"channel": {
			"status": status,
			"game":   game,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", u, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Set("Authorization", "OAuth "+token)
	req.Header.Set("Client-ID", t.clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code %d", resp.StatusCode)
	}

	return nil
}

// Game represents information about a game on Twitch.
type Game struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Image      string `json:"image"`
}

type byPopularity []Game

func (g byPopularity) Len() int           { return len(g) }
func (g byPopularity) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g byPopularity) Less(i, j int) bool { return g[i].Popularity > g[j].Popularity }

// Games returns what games are available.
func (t *API) Games() []Game {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.games == nil {
		// TODO: consider polling periodically to refresh cache
		t.discoverGames()
	}
	return t.games
}

func (t *API) discoverGames() {
	out := make([]Game, 0, 2048)

	for offset, total := 0, 100; offset < total; offset += 100 {
		var (
			g   []Game
			err error
		)
		g, total, err = t.makeGamesRequest(offset)
		if err != nil {
			log.Printf("got err while making games request: %s", err)
			continue
		}
		out = append(out, g...)
	}

	sort.Sort(byPopularity(out))
	result := make([]Game, 0, len(out))
	idSet := make(map[int]struct{})
	for _, g := range out {
		if _, ok := idSet[g.ID]; ok {
			continue
		}
		idSet[g.ID] = struct{}{}
		result = append(result, g)
	}
	t.games = result
}

func (t *API) makeGamesRequest(offset int) ([]Game, int, error) {
	req, err := t.buildGamesRequest(offset)
	if err != nil {
		return nil, 0, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("Bad status code %d", resp.StatusCode)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("err closing response body: %s", err)
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var r map[string]interface{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, 0, err
	}

	ftotal, ok := r["_total"].(float64)
	if !ok {
		return nil, 0, errors.New("unable to type assert _total")
	}
	total := int(ftotal)
	top, ok := r["top"].([]interface{})
	if !ok {
		return nil, 0, errors.New("unable to type assert top")
	}
	var games []Game
	for _, gp := range top {
		gpayload, ok := gp.(map[string]interface{})
		if !ok {
			log.Println("unable to type assert game payload")
			continue
		}
		g, ok := gpayload["game"].(map[string]interface{})
		if !ok {
			log.Println("unable to type assert game")
			continue
		}
		name, ok := g["name"].(string)
		if !ok {
			log.Println("unable to type assert name")
			continue
		}
		fpopularity, ok := g["popularity"].(float64)
		if !ok {
			log.Println("unable to type assert popularity")
			continue
		}
		popularity := int(fpopularity)
		fid, ok := g["_id"].(float64)
		if !ok {
			log.Println("unable to type assert _id")
			continue
		}
		id := int(fid)
		box, ok := g["box"].(map[string]interface{})
		if !ok {
			log.Println("unable to type assert box")
			continue
		}
		image, ok := box["small"].(string)
		if !ok {
			log.Println("unable to type assert small")
			continue
		}
		games = append(games, Game{
			ID:         id,
			Name:       name,
			Popularity: popularity,
			Image:      image,
		})
	}
	return games, total, nil
}

func (t *API) buildGamesRequest(offset int) (*http.Request, error) {
	req, err := http.NewRequest("GET", t.url+"/games/top", nil)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("limit", "100")
	v.Set("offset", strconv.Itoa(offset))
	req.URL.RawQuery = v.Encode()

	req.Header.Set("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Set("Client-ID", t.clientID)

	return req, nil
}
