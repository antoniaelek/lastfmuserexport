package lastfmexport

import (
	"sort"
	"strconv"
	"time"
)

// Scrobble is a scrobbled track.
type Scrobble struct {
	Track     string
	Artist    string
	Album     string
	Timestamp time.Time
	URL       string
}

type recentTracksResponse struct {
	Recenttracks struct {
		Attr struct {
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			User       string `json:"user"`
			Total      string `json:"total"`
			TotalPages string `json:"totalPages"`
		} `json:"@attr"`
		Track []struct {
			Artist struct {
				Mbid string `json:"mbid"`
				Text string `json:"#text"`
			} `json:"artist"`
			Album struct {
				Mbid string `json:"mbid"`
				Text string `json:"#text"`
			} `json:"album"`
			Image []struct {
				Size string `json:"size"`
				Text string `json:"#text"`
			} `json:"image"`
			Streamable string `json:"streamable"`
			Date       struct {
				Uts  string `json:"uts"`
				Text string `json:"#text"`
			} `json:"date"`
			URL  string `json:"url"`
			Name string `json:"name"`
			Mbid string `json:"mbid"`
		} `json:"track"`
	} `json:"recenttracks"`
}

// GetScrobbled gets user's scrobbled tracks.
func GetScrobbled(username string, apiKey string) (tracks []Scrobble, err error) {
	resp := new(recentTracksResponse)
	getJSON("http://ws.audioscrobbler.com/2.0/?"+
		"method=user.getrecenttracks"+
		"&api_key="+apiKey+
		"&format=json"+
		"&user="+username+
		"&page=1", resp)

	total, err := strconv.Atoi(resp.Recenttracks.Attr.Total)
	if err != nil {
		return
	}

	totalPages, err := strconv.Atoi(resp.Recenttracks.Attr.TotalPages)
	if err != nil {
		return
	}

	messages := make(chan *recentTracksResponse)

	for i := 1; i <= totalPages; i++ {
		go getPage(i, messages, username, apiKey)
	}

	tracks = make([]Scrobble, total)

	// Writing down the way the standard time would look like formatted our way
	// Standard time is "Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"
	layout := "02 Jan 2006, 15:04"

	for i := 0; i < totalPages; i++ {
		resp := <-messages
		ts := resp.Recenttracks.Track
		for j, track := range ts {
			scrobbleTime, _ := time.Parse(layout, track.Date.Text)
			t := Scrobble{
				Track:     track.Name,
				Artist:    track.Artist.Text,
				Album:     track.Album.Text,
				Timestamp: scrobbleTime,
				URL:       track.URL,
			}
			tracks[i*50+j] = t
		}
	}

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Timestamp.Before(tracks[j].Timestamp)
	})

	return tracks, nil
}

func getPage(page int, c chan *recentTracksResponse, username string, apiKey string) {
	resp := new(recentTracksResponse)
	for {
		getJSON(baseURL+
			"?method=user.getrecenttracks"+
			"&api_key="+apiKey+
			"&format=json"+
			"&user="+username+
			"&page="+strconv.Itoa(page), resp)
		if len(resp.Recenttracks.Track) > 0 {
			c <- resp
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}
