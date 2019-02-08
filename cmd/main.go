package main

import (
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

	// Get scrobbles
	start := time.Now()
	tracks, err := lastfmexport.GetScrobbled("muser1901", config.APIKey)
	elapsed := time.Since(start)
	log.Printf("Execution time: %s\n", elapsed)

	// Print scrobbles
	if err != nil {
		log.Fatalln("Error fetching scrobbles:", err)
	} else {
		for _, t := range tracks {
			log.Printf("%s\t%s\t%s\t%s\n", t.Timestamp, t.Track, t.Artist, t.URL)
		}
	}
}
