package export

import (
	"encoding/json"
	"net/http"
)

const baseURL = "http://ws.audioscrobbler.com/2.0/"

func getJSON(url string, client *http.Client, target interface{}) error {
	for {
		r, err := client.Get(url)
		if err == nil && r.StatusCode == 200 {
			defer r.Body.Close()
			return json.NewDecoder(r.Body).Decode(target)
		}
	}
}
