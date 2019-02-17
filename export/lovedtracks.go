package export

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Track is a track from LastFm
type Track struct {
	Timestamp time.Time
	Name      string
	Artist    string
	URL       string
}

// TrackArray is an array of Track objects
type TrackArray []Track

// ToCsv converts array of Track objects to csv
func (Tracks TrackArray) ToCsv(sep string) []string {
	csv := make([]string, len(Tracks))
	for i, Track := range Tracks {
		csv[i] = Track.Timestamp.String() + sep + Track.Name + sep + Track.Artist + sep + Track.URL
	}
	return csv
}

type lovedTracksResponse struct {
	Lovedtracks struct {
		Attr struct {
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			User       string `json:"user"`
			Total      string `json:"total"`
			TotalPages string `json:"totalPages"`
		} `json:"@attr"`
		Track []struct {
			Artist struct {
				URL  string `json:"url"`
				Name string `json:"name"`
				Mbid string `json:"mbid"`
			} `json:"artist"`
			Mbid string `json:"mbid"`
			Date struct {
				Uts  string `json:"uts"`
				Text string `json:"#text"`
			} `json:"date"`
			URL   string `json:"url"`
			Image []struct {
				Size string `json:"size"`
				Text string `json:"#text"`
			} `json:"image"`
			Name       string `json:"name"`
			Streamable struct {
				Fulltrack string `json:"fulltrack"`
				Text      string `json:"#text"`
			} `json:"streamable"`
		} `json:"track"`
	} `json:"lovedtracks"`
}

// GetLovedTracks gets user's loved tracks
func GetLovedTracks(username string, apiKey string) (tracks []Track, err error) {
	var client = http.Client{Timeout: 10 * time.Second}
	resp := new(lovedTracksResponse)
	getJSON(baseURL+
		"?method=user.getlovedtracks"+
		"&user="+username+
		"&api_key="+apiKey+
		"&format=json", &client, resp)

	total, err := strconv.Atoi(resp.Lovedtracks.Attr.Total)
	if err != nil {
		return
	}
	total, err = strconv.Atoi(resp.Lovedtracks.Attr.Total)
	if err != nil {
		return
	}

	totalPages, err := strconv.Atoi(resp.Lovedtracks.Attr.TotalPages)
	if err != nil {
		return
	}

	log.Printf("There are %d loved tracks across %d pages\n", total, totalPages)

	chunkSize := 30
	tracks = make([]Track, 0, total)
	for i := 1; i <= totalPages; i = i + chunkSize {
		upperBound := i + chunkSize - 1
		if upperBound > totalPages {
			upperBound = totalPages
		}
		tracks = append(tracks, getlovedTracksPart(&client, i, upperBound, username, apiKey)...)
	}

	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Timestamp.Before(tracks[j].Timestamp)
	})

	return tracks, nil
}

func getlovedTracksPart(client *http.Client, firstPage int, lastPage int, username string, apiKey string) []Track {
	messages := make(chan *lovedTracksResponse)

	for i := firstPage; i <= lastPage; i++ {
		go getLovedTracksPage(i, client, messages, username, apiKey)
	}

	// Writing down the way the standard time would look like formatted our way
	// Standard time is "Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"
	layout := "02 Jan 2006, 15:04"

	var tracks []Track
	for i := firstPage; i <= lastPage; i++ {
		resp := <-messages
		ts := resp.Lovedtracks.Track
		for _, track := range ts {
			scrobbleTime, _ := time.Parse(layout, track.Date.Text)
			t := Track{
				Name:      track.Name,
				Artist:    track.Artist.Name,
				Timestamp: scrobbleTime,
				URL:       track.URL,
			}
			tracks = append(tracks, t)
		}
	}
	return tracks
}

func getLovedTracksPage(page int, client *http.Client, c chan *lovedTracksResponse, username string, apiKey string) {
	resp := new(lovedTracksResponse)
	for {
		getJSON(baseURL+
			"?method=user.getlovedtracks"+
			"&user="+username+
			"&api_key="+apiKey+
			"&format=json"+
			"&page="+strconv.Itoa(page), client, resp)

		if len(resp.Lovedtracks.Track) > 0 {
			log.Printf("OK loved tracks page %d\n", page)
			c <- resp
			break
		}

		log.Printf("RETRY loved tracks page %d\n", page)
		time.Sleep(500 * time.Millisecond)
	}
}
