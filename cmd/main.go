package main

import (
	"fmt"
	"lastfmuserexport"
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

	// Get
	data, err := lastfmuserexport.GetTags("muser1901", config.APIKey)

	// Print
	if err != nil {
		log.Fatalln("Error fetching data:", err)
	} else {
		var x lastfmuserexport.TagArray
		x = data
		csv := x.ToCsv("\t")
		for _, l := range csv {
			fmt.Println(l)
		}
	}

	// Print exec time
	elapsed := time.Since(start)
	log.Printf("Execution time: %s\n", elapsed)
}
