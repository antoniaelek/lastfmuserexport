package lastfmuserexport

import (
	"encoding/json"
	"net/http"
	"time"
)

const baseURL = "http://ws.audioscrobbler.com/2.0/"

func getJSON(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	for {
		r, err := myClient.Get(url)
		if err == nil && r.StatusCode == 200 {
			defer r.Body.Close()
			return json.NewDecoder(r.Body).Decode(target)
		}
	}
}
