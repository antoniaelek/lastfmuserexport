package main

import (
	"fmt"
	"lastfmexport"
	"log"
	"time"

	"github.com/tkanos/gonfig"
)

// Config represents application configuration
type Config struct {
	APIKey string
}

func main() {

	// Get config
	config := Config{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatalln("Error fetching config:", err)
		return
	}

	// Measure exec time
	start := time.Now()

	// Get scrobbles
	tracks, err := lastfmexport.GetLovedTracks("muser1901", config.APIKey)

	// Print scrobbles
	if err != nil {
		log.Fatalln("Error fetching scrobbles:", err)
	} else {
		var x lastfmexport.TrackArray
		x = tracks
		csv := x.ToCsv("\t")
		for _, l := range csv {
			fmt.Println(l)
		}
	}

	// Print exec time
	elapsed := time.Since(start)
	log.Printf("Execution time: %s\n", elapsed)
}
