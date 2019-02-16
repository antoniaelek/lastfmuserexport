package export

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Tag is a tag
type Tag struct {
	Count int
	Name  string
	URL   string
}

// TagArray is an array of Tag objects
type TagArray []Tag

// ToCsv converts array of Tag objects to csv
func (Tags TagArray) ToCsv(sep string) []string {
	csv := make([]string, len(Tags))
	for i, Tag := range Tags {
		csv[i] = strconv.Itoa(Tag.Count) + sep + Tag.Name + sep + Tag.URL
	}
	return csv
}

type tagsResponse struct {
	Toptags struct {
		Tag []struct {
			Count int    `json:"count"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"tag"`
		Attr struct {
			Artist string `json:"artist"`
		} `json:"@attr"`
	} `json:"toptags"`
}

// GetTags gets top tags for user
func GetTags(user string, apiKey string) (tags []Tag, err error) {
	artists, err := GetArtists(user, apiKey)
	if err != nil {
		return nil, err
	}

	tagsMap := make(map[string]Tag)
	for _, artist := range artists {
		artistTags, err := getTagsForArtist(artist.Name, apiKey)
		if err != nil {
			fmt.Println("ERROR")
			continue
		}
		for _, at := range artistTags {
			tag, exists := tagsMap[at.Name]
			if exists {
				at.Count += tag.Count
			}
			tagsMap[at.Name] = at
		}
	}

	tags = make([]Tag, len(tagsMap))
	idx := 0
	for _, tag := range tagsMap {
		tags[idx] = tag
		idx++
	}

	return tags, nil
}

func getTagsForArtist(artist string, apiKey string) (tags []Tag, err error) {
	resp := new(tagsResponse)
	getJSON(baseURL+
		"?method=artist.gettoptags"+
		"&api_key="+apiKey+
		"&format=json"+
		"&artist="+strings.Replace(artist, " ", "+", -1), resp)

	tags = make([]Tag, len(resp.Toptags.Tag))
	for i, t := range resp.Toptags.Tag {
		tag := Tag{
			Name:  t.Name,
			Count: t.Count,
			URL:   t.URL,
		}
		tags[i] = tag
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	return tags, nil
}
