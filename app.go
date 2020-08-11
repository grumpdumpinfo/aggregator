package main

import (
	"encoding/json"
	"fmt"
	"github.com/grumpdumpinfo/aggregator/data"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type App struct {
	TargetPlaylists []string
	APIKey string
	L *zap.Logger
}

// GetAllPlaylistsFirstPage gets the first 20 items of each playlist in TargetPlaylists
func (a *App) GetAllPlaylistsFirstPage() []data.Playlist {
	result := make([]data.Playlist, 0, len(a.TargetPlaylists))

	for _, playlist := range a.TargetPlaylists {
		maybePlaylist := a.GetPlaylistInfoPage(playlist, "")
		if maybePlaylist != nil {
			result = append(result, *maybePlaylist)
		}
	}

	return result
}

// GetAllPlaylistsPages gets every page of every playlist in App
// In the event of a failed fetch, GetAllPlaylistsPages will return an array of what pages it _could_ retrieve
func (a *App) GetAllPlaylistsPages() []data.Playlist {
	result := make([]data.Playlist, 0, len(a.TargetPlaylists))

	for _, playlist := range a.TargetPlaylists {
		maybePlaylist := a.GetAllPagesForPlaylist(playlist)
		if maybePlaylist != nil {
			result = append(result, *maybePlaylist)
		}
	}

	return result
}

// GetAllPagesForPlaylist gets every page for a given playlist
//
func (a *App) GetAllPagesForPlaylist(playlist string) *data.Playlist {
	firstPage := a.GetPlaylistInfoPage(playlist, "")

	if firstPage == nil {
		a.L.Error("could not get first page of playlist",
			zap.String("playlist", playlist),
		)

		return nil
	}

	nextPageToken := firstPage.NextPageToken

	for {
		nextPage := a.GetPlaylistInfoPage(playlist, nextPageToken)

		if nextPage == nil {
			a.L.Error("expected nextPage in playlist, got nil",
				zap.String("playlist", playlist),
				zap.String("nextPageToken", nextPageToken),
			)
			// We need to end iteration when nextPage is nil, since we can't skip pages
			break
		}

		nextPageToken = nextPage.NextPageToken
		// This might be a performance bottleneck, if it becomes a problem change this
		firstPage.Items = append(firstPage.Items, nextPage.Items...)

		if nextPageToken == "" {
			break
		}
	}

	return firstPage
}

// GetPlaylistInfoPage gets a given page of playlist information for the application's playlist
// if page is an empty string, returns the first page of the request
func (a *App) GetPlaylistInfoPage(playlist, page string) *data.Playlist {
	var req string
	if page == "" {
		req = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=20&playlistId=%s&key=%s&Accept=application/json",
			playlist,
			a.APIKey)
	} else {
		req = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=20&playlistId=%s&key=%s&pageToken=%s&Accept=application/json",
			playlist,
			a.APIKey,
			page)
	}

	start := time.Now()

	a.L.Info("beginning download of playlist information",
		zap.Time("now", time.Now()),
		zap.String("page", page),
		zap.String("playlist", playlist),
	)

	// Seriously, for the love of god, do not let this touch prod without doing validation first
	res, err := http.Get(req)

	if err != nil {
		a.L.Error("failed to download playlist information",
			zap.String("page", page),
			zap.String("playlist", playlist),
			zap.Error(err),
		)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)

	a.L.Info("finished downloading playlist information",
		zap.Time("now", time.Now()),
		zap.Duration("roundtrip", time.Since(start)),
		zap.String("page", page),
		zap.String("playlist", playlist),
	)

	if err != nil {
		a.L.Error("failed to read playlist information to byte array",
			zap.String("page", page),
			zap.String("playlist", playlist),
			zap.Error(err),
		)

		return nil
	}

	result := new(data.Playlist)
	err = json.Unmarshal(body, result)

	if err != nil {
		a.L.Error("failed to unmarshal response JSON",
			zap.String("page", page),
			zap.String("playlist", playlist),
			zap.Error(err),
		)

		return nil
	}

	return result
}