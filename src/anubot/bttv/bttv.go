package bttv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aymerick/raymond"
)

const baseURL = "https://api.betterttv.net/2/"

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}

type apiResponse struct {
	Emotes []struct {
		ID   string `json:"id"`
		Code string `json:"code"`
	} `json:"emotes"`
	URLTemplate string `json:"urlTemplate"`
}

// Emoji returns the current BTTV emoji. Optionally, you may provide a channel
// name and it will get the emoji specific for that channel.
func Emoji(channel string) (map[string]string, error) {
	result := make(map[string]string)

	err := requestEmoji("emotes", result)
	if err != nil {
		return nil, err
	}

	if channel != "" {
		_ = requestEmoji("channels/"+channel, result)
	}

	return result, nil
}

func requestEmoji(path string, result map[string]string) error {
	resp, err := httpClient.Get(baseURL + path)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got bad status code from bttv api: %d", resp.StatusCode)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("got err while closing resp body: %s", err)
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var decoded apiResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		return err
	}

	for _, e := range decoded.Emotes {
		if e.Code == "" {
			log.Print("empty code for bttv emote")
			continue
		}
		rendered, err := raymond.Render("https:"+decoded.URLTemplate, map[string]string{
			"id":    e.ID,
			"image": "1x",
		})
		if err != nil {
			log.Printf("got error while rendering bttv url template: %s", err)
			continue
		}
		_, err = url.Parse(rendered)
		if err != nil {
			log.Printf("got error while parsing bttv url template: %s", err)
			continue
		}
		result[e.Code] = rendered
	}

	return nil
}
