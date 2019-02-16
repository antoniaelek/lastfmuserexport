package export

import (
	"log"
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
	resp := new(artistsResponse)
	getJSON(baseURL+
		"?method=user.gettopartists"+
		"&api_key="+apiKey+
		"&format=json"+
		"&user="+user, resp)

	total, err := strconv.Atoi(resp.Topartists.Attr.Total)
	if err != nil {
		return
	}

	totalPages, err := strconv.Atoi(resp.Topartists.Attr.TotalPages)
	if err != nil {
		return
	}

	messages := make(chan *artistsResponse)

	for i := 1; i <= totalPages; i++ {
		go getArtistsPage(i, messages, user, apiKey)

		// Because of rate limiting
		if i%20 == 0 {
			time.Sleep(1000 * time.Millisecond)
		}
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

func getArtistsPage(page int, c chan *artistsResponse, username string, apiKey string) {
	resp := new(artistsResponse)
	for {
		getJSON(baseURL+
			"?method=user.gettopartists"+
			"&user="+username+
			"&page="+strconv.Itoa(page)+
			"&api_key="+apiKey+
			"&format=json", resp)

		if len(resp.Topartists.Artist) > 0 {
			log.Printf("%-5s artists page %d\n", "OK", page)
			c <- resp
			break
		}

		log.Printf("%-5s artists page %d\n", "RETRY", page)
		time.Sleep(500 * time.Millisecond)
	}
}
