package export

import (
	"log"
	"net/http"
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

// ScrobbleArray is an array of scrobble objects
type ScrobbleArray []Scrobble

// ToCsv converts array of scrobble objects to csv
func (scrobbles ScrobbleArray) ToCsv(sep string) []string {
	csv := make([]string, len(scrobbles))
	for i, scrobble := range scrobbles {
		csv[i] = scrobble.Timestamp.String() + sep + scrobble.Track + sep + scrobble.Artist + sep + scrobble.Album + sep + scrobble.URL
	}
	return csv
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

// GetScrobbles gets user's scrobbled tracks.
func GetScrobbles(username string, apiKey string) (tracks []Scrobble, err error) {
	var client = http.Client{Timeout: 10 * time.Second}

	resp := new(recentTracksResponse)
	getJSON(baseURL+
		"?method=user.getrecenttracks"+
		"&api_key="+apiKey+
		"&format=json"+
		"&user="+username+
		"&page=1", &client, resp)

	total, err := strconv.Atoi(resp.Recenttracks.Attr.Total)
	if err != nil {
		return
	}

	totalPages, err := strconv.Atoi(resp.Recenttracks.Attr.TotalPages)
	if err != nil {
		return
	}

	log.Printf("There are %d scrobbles across %d pages\n", total, totalPages)

	chunkSize := 30
	tracks = make([]Scrobble, 0, total)
	for i := 1; i <= totalPages; i = i + chunkSize {
		upperBound := i + chunkSize - 1
		if upperBound > totalPages {
			upperBound = totalPages
		}
		tracks = append(tracks, getPart(&client, i, upperBound, username, apiKey)...)
	}

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Timestamp.Before(tracks[j].Timestamp)
	})

	return tracks, nil
}

func getPart(client *http.Client, firstPage int, lastPage int, username string, apiKey string) []Scrobble {
	messages := make(chan *recentTracksResponse)

	for i := firstPage; i <= lastPage; i++ {
		go getPage(i, client, messages, username, apiKey)
	}

	// Writing down the way the standard time would look like formatted our way
	// Standard time is "Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"
	layout := "02 Jan 2006, 15:04"

	var tracks []Scrobble
	for i := firstPage; i <= lastPage; i++ {
		resp := <-messages
		ts := resp.Recenttracks.Track
		for _, track := range ts {
			scrobbleTime, _ := time.Parse(layout, track.Date.Text)
			t := Scrobble{
				Track:     track.Name,
				Artist:    track.Artist.Text,
				Album:     track.Album.Text,
				Timestamp: scrobbleTime,
				URL:       track.URL,
			}
			tracks = append(tracks, t)
		}
	}
	return tracks
}

func getPage(page int, client *http.Client, c chan *recentTracksResponse, username string, apiKey string) {
	resp := new(recentTracksResponse)
	i := 1
	for {
		getJSON(baseURL+
			"?method=user.getrecenttracks"+
			"&api_key="+apiKey+
			"&format=json"+
			"&user="+username+
			"&page="+strconv.Itoa(page), client, resp)

		if len(resp.Recenttracks.Track) > 0 {
			log.Printf("OK scrobbles page %d\n", page)
			c <- resp
			break
		}

		log.Printf("RETRY %-2d scrobbles page %d\n", i, page)
		i++
	}
}
