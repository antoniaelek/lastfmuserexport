package export

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Artist is scrobbled artist
type Artist struct {
	Name      string
	URL       string
	PlayCount int
}

// ArtistArray is an array of Artist objects
type ArtistArray []Artist

// ToCsv converts array of Artist objects to csv
func (Artists ArtistArray) ToCsv(sep string) []string {
	csv := make([]string, len(Artists))
	for i, Artist := range Artists {
		csv[i] = strconv.Itoa(Artist.PlayCount) + sep + Artist.Name + sep + Artist.URL
	}
	return csv
}

type artistsResponse struct {
	Topartists struct {
		Artist []struct {
			Attr struct {
				Rank string `json:"rank"`
			} `json:"@attr"`
			Mbid      string `json:"mbid"`
			URL       string `json:"url"`
			Playcount string `json:"playcount"`
			Image     []struct {
				Size string `json:"size"`
				Text string `json:"#text"`
			} `json:"image"`
			Name       string `json:"name"`
			Streamable string `json:"streamable"`
		} `json:"artist"`
		Attr struct {
			Page       string `json:"page"`
			PerPage    string `json:"perPage"`
			User       string `json:"user"`
			Total      string `json:"total"`
			TotalPages string `json:"totalPages"`
		} `json:"@attr"`
	} `json:"topartists"`
}

// GetArtists gets top tags for an artist
func GetArtists(user string, apiKey string) (artists []Artist, err error) {
	var client = &http.Client{Timeout: 10 * time.Second}
	resp := new(artistsResponse)
	getJSON(baseURL+
		"?method=user.gettopartists"+
		"&api_key="+apiKey+
		"&format=json"+
		"&user="+user, client, resp)

	total, err := strconv.Atoi(resp.Topartists.Attr.Total)
	if err != nil {
		return
	}

	totalPages, err := strconv.Atoi(resp.Topartists.Attr.TotalPages)
	if err != nil {
		return
	}

	log.Printf("There are %d artistss across %d pages\n", total, totalPages)

	messages := make(chan *artistsResponse)

	for i := 1; i <= totalPages; i++ {
		go getArtistsPage(i, client, messages, user, apiKey)
	}

	artists = make([]Artist, total)
	idx := 0
	for i := 0; i < totalPages; i++ {
		resp := <-messages
		ts := resp.Topartists.Artist
		for _, a := range ts {
			cnt, err := strconv.Atoi(a.Playcount)
			if err != nil {
				continue
			}
			artist := Artist{
				Name:      a.Name,
				PlayCount: cnt,
				URL:       a.URL,
			}
			artists[idx] = artist
			idx++
		}
	}

	sort.Slice(artists, func(i, j int) bool {
		return artists[i].Name < artists[j].Name
	})

	return artists, nil
}

func getArtistsPage(page int, client *http.Client, c chan *artistsResponse, username string, apiKey string) {
	resp := new(artistsResponse)
	for {
		getJSON(baseURL+
			"?method=user.gettopartists"+
			"&user="+username+
			"&page="+strconv.Itoa(page)+
			"&api_key="+apiKey+
			"&format=json", client, resp)

		if len(resp.Topartists.Artist) > 0 {
			log.Printf("OK artists page %d\n", page)
			c <- resp
			break
		}

		log.Printf("RETRY artists page %d\n", page)
		time.Sleep(500 * time.Millisecond)
	}
}
