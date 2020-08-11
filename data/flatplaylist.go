package data

import "encoding/json"

// FlattenedPlaylist essentially represents the same type of data as Playlist, but is in a format that's easier to
// deal with.
type FlattenedPlaylist struct {
	PublishedAt string `json:"publishedAt"`
	Title string `json:"title"`
	Description string `json:"description"`
	Thumbnails json.RawMessage `json:"thumbnails"`
	VideoID string `json:"videoId"`
	ChannelTitle string `json:"channelTitle"`
}

func FlattenedPlaylistsFromPlaylists(playlists []Playlist) []FlattenedPlaylist {
	count := 0

	for _, playlist := range playlists {
		count += len(playlist.Items)
	}

	toReturn := make([]FlattenedPlaylist, 0, count)

	for _, playlist := range playlists {
		toReturn = append(toReturn, FlattenedPlaylistsFromPlaylist(playlist)...)
	}

	return toReturn
}

func FlattenedPlaylistsFromPlaylist(playlist Playlist) []FlattenedPlaylist {
	toReturn := make([]FlattenedPlaylist, 0, len(playlist.Items))

	for _, element := range playlist.Items {
		toAdd := FlattenedPlaylist{
			PublishedAt: element.Snippet.PublishedAt,
			Title: element.Snippet.Title,
			Description: element.Snippet.Description,
			Thumbnails: element.Snippet.Thumbnails,
			// this right here is the only reason FlattenedPlaylist isn't embedded in playlist
			// consider changing this
			VideoID: element.Snippet.ResourceID.VideoID,
			ChannelTitle: element.Snippet.ChannelTitle,
		}

		toReturn = append(toReturn, toAdd)
	}

	return toReturn
}