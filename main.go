package main

import (
	"flag"
	"fmt"
	"lastfmuserexport/export"
	"log"
	"os"
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
			log.Println(e)
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
	log.Printf("Fetching scrobbles for user %s with api key %s...\n", user, apiKey)
	data, err := export.GetScrobbles(user, apiKey)

	if err != nil {
		log.Println("Error fetching scrobbles.")
		panic(err)
	}

	log.Println("Fetched scrobbles.")
	var x export.ScrobbleArray
	x = data
	csv := x.ToCsv("\t")
	file := "scrobbles.csv"

	err = SaveCsvFile(csv, file)
	if err != nil {
		log.Panicf("Error saving to %s.\n", file)
		panic(err)
	}
	log.Printf("Saved %s.\n", file)
}

// GetLovedTracks gets loved tracks
func GetLovedTracks(user string, apiKey string) {
	log.Printf("Fetching loved tracks for user %s with api key %s...\n", user, apiKey)
	data, err := export.GetLovedTracks(user, apiKey)

	if err != nil {
		log.Println("Error fetching loved tracks:")
		panic(err)
	}

	var x export.TrackArray
	x = data
	csv := x.ToCsv("\t")
	file := "loved.csv"

	log.Println("Fetched loved tracks.")
	err = SaveCsvFile(csv, file)
	if err != nil {
		log.Panicf("Error saving to %s.\n", file)
		panic(err)
	}

	log.Printf("Saved %s.\n", file)
}

// GetTags gets tags
func GetTags(user string, apiKey string) {
	log.Printf("Fetching tags for user %s with api key %s...\n", user, apiKey)

	data, err := export.GetTags(user, apiKey)
	if err != nil {
		log.Println("Error fetching tags.")
		panic(err)
	}

	log.Println("Fetched tags.")
	var x export.TagArray
	x = data
	csv := x.ToCsv("\t")
	file := "tags.csv"

	err = SaveCsvFile(csv, file)
	if err != nil {
		log.Panicf("Error saving to %s.\n", file)
		panic(err)
	}

	log.Printf("Saved %s.\n", file)
}

// GetArtists gets artists
func GetArtists(user string, apiKey string) {
	log.Printf("Fetching artists for user %s with api key %s...\n", user, apiKey)

	data, err := export.GetArtists(user, apiKey)

	if err != nil {
		log.Println("Error fetching artists.")
		panic(err)
	}

	log.Println("Fetched artists.")
	var x export.ArtistArray
	x = data
	csv := x.ToCsv("\t")
	file := "artists.csv"
	err = SaveCsvFile(csv, file)
	if err != nil {
		log.Panicf("Error saving to %s.\n", file)
		panic(err)
	}

	log.Printf("Saved %s.\n", file)
}

// SaveCsvFile saves csv to file
func SaveCsvFile(csv []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range csv {
		fmt.Fprintln(f, line)
	}
	defer f.Close()

	return nil
}
