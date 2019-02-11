package lastfmexport

// Tag is a tag
type Tag struct {
	Count int
	Name  string
	URL   string
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

// func GetTopTagsForUser(user string, apiKey string) (tags []Tag, err error) {
// }

// GetTopTags gets top tags for an artist
func GetTopTagsForArtist(artist string, apiKey string) (tags []Tag, err error) {
	resp := new(tagsResponse)
	getJSON("http://ws.audioscrobbler.com/2.0/?"+
		"method=artist.gettoptags"+
		"&api_key="+apiKey+
		"&format=json"+
		"&artist="+artist, resp)

	tags = make([]Tag, len(resp.Toptags.Tag))
	for i, t := range resp.Toptags.Tag {
		tag := Tag{
			Name:  t.Name,
			Count: t.Count,
			URL:   t.URL,
		}
		tags[i] = tag
	}

	return tags, nil
}
