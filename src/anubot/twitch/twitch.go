package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

// New creates a new API.
func New(url, clientID string) API {
	if url == "" {
		url = twitchAPIURL
	}
	return API{
		url:      url,
		clientID: clientID,
	}
}

// Username gets the username for a give oauth token.
func (t API) Username(token string) (username string, err error) {
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

	defer resp.Body.Close()
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
func (t API) StreamInfo(channel string) (string, string, error) {
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

	defer resp.Body.Close()
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
func (t API) UpdateDescription(status, game, channel, token string) error {
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
