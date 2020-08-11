package data

import "encoding/json"

type Playlist struct {
	NextPageToken string `json:"nextPageToken"`
	PageInfo struct{
		TotalResults uint `json:"totalResults"`
		ResultsPerPage uint `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		ID string `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelID string `json:"channelId"`
			Title string `json:"title"`
			Description string `json:"description"`
			// the thumbnails is a JSON structure too, but I plan to pass it to the frontend unmodified, so parsing it
			// is a waste of time... probably.
			Thumbnails json.RawMessage `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
			ResourceID struct {
				VideoID string `json:"videoId"`
			} `json:"resouceId"`
		} `json:"snippet"`
	} `json:"items"`
}