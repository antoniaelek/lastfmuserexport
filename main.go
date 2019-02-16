package main

import (
	"bufio"
	"flag"
	"fmt"
	"lastfmuserexport/export"
	"log"
	"os"
	"time"
)

func main() {
	key := flag.String("key", "", "Last.fm API key (required)")
	user := flag.String("user", "", "Last.fm username (required)")

	scrobbles := flag.Bool("scrobbles", false, "Save scrobbles to 'scrobbles.csv' file")
	loved := flag.Bool("loved", false, "Save loved tracks to 'loved.csv' file")
	artists := flag.Bool("artists", false, "Save artists to 'artists.csv' file")
	tags := flag.Bool("tags", false, "Save tags to 'tags.csv' file")

	flag.Parse()

	var err []string
	if *key == "" {
		err = append(err, "Missing required parameter: key")
	}
	if *user == "" {
		err = append(err, "Missing required parameter: user")
	}

	if err != nil && len(err) > 0 {
		for _, e := range err {
			fmt.Println(e)
		}
		return
	}

	if *scrobbles {
		GetSrobbles(*user, *key)
	}
	if *loved {
		GetLovedTracks(*user, *key)
	}
	if *artists {
		GetArtists(*user, *key)
	}
	if *tags {
		GetTags(*user, *key)
	}
}

// GetSrobbles gets scrobbles
func GetSrobbles(user string, apiKey string) {
	fmt.Println("Fetching scrobbles...")
	data, err := export.GetScrobbles(user, apiKey)

	if err != nil {
		fmt.Println("Error fetching scrobbles:", err)
	} else {
		var x export.ScrobbleArray
		x = data
		csv := x.ToCsv("\t")
		err = SaveCsvFile(csv, "scrobbles.csv")
		if err != nil {
			fmt.Println("Error saving scrobbles file", err)
		}
	}
}

// GetLovedTracks gets loved tracks
func GetLovedTracks(user string, apiKey string) {
	fmt.Println("Fetching loved tracks...")
	data, err := export.GetLovedTracks(user, apiKey)

	if err != nil {
		fmt.Println("Error fetching loved tracks:", err)
	} else {
		var x export.TrackArray
		x = data
		csv := x.ToCsv("\t")
		err = SaveCsvFile(csv, "loved.csv")
		if err != nil {
			fmt.Println("Error saving loved tracks file", err)
		}
	}
}

// GetTags gets tags
func GetTags(user string, apiKey string) {
	fmt.Println("Fetching tags...")
	data, err := export.GetTags(user, apiKey)

	if err != nil {
		fmt.Println("Error fetching tags:", err)
	} else {
		var x export.TagArray
		x = data
		csv := x.ToCsv("\t")
		err = SaveCsvFile(csv, "tags.csv")
		if err != nil {
			fmt.Println("Error saving tags file", err)
		}
	}
}

// GetArtists gets artists
func GetArtists(user string, apiKey string) {
	fmt.Println("Fetching artists...")
	start := time.Now()
	data, err := export.GetArtists(user, apiKey)
	elapsed := time.Since(start)
	log.Printf("Artists fetch time: %s\n", elapsed)

	if err != nil {
		fmt.Println("Error fetching artists:", err)
	} else {
		var x export.ArtistArray
		x = data
		csv := x.ToCsv("\t")
		err = SaveCsvFile(csv, "artists.csv")
		if err != nil {
			fmt.Println("Error saving artists file", err)
		}
	}
}

// SaveCsvFile saves csv to file
func SaveCsvFile(csv []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, line := range csv {
		w.WriteString(line + "\n")
	}
	return nil
}
