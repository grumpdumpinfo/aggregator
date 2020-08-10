package main

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type App struct {
	TargetPlaylist string
	APIKey string
	L *zap.Logger
}

// GetPlaylistInfoPage gets a given page of playlist information for the application's playlist
// if page is an empty string, returns the first page of the request
func (a *App) GetPlaylistInfoPage(page string) []byte {
	// TODO: validate or escape input to protect from injection attacks
	// Validation should also be done inside NewApp, however we still shouldn't assume that the user will use app
	// correctly.
	var req string
	if page == "" {
		req = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=20&playlistId=%s&key=%s&Accept=application/json",
			a.TargetPlaylist,
			a.APIKey)
	} else {
		req = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=20&playlistId=%s&key=%s&pageToken=%s&Accept=application/json",
			a.TargetPlaylist,
			a.APIKey,
			page)
	}

	start := time.Now()

	a.L.Info("beginning download of playlist information",
		zap.Time("now", time.Now()),
		zap.String("page", page),
		zap.String("playlist", a.TargetPlaylist),
	)

	// Seriously, for the love of god, do not let this touch prod without doing validation first
	res, err := http.Get(req)

	if err != nil {
		a.L.Error("failed to download playlist information",
			zap.String("page", page),
			zap.String("playlist", a.TargetPlaylist),
			zap.Error(err),
		)
		return nil
	}

	toReturn, err := ioutil.ReadAll(res.Body)

	a.L.Info("finished downloading playlist information",
		zap.Time("now", time.Now()),
		zap.Duration("roundtrip", time.Since(start)),
		zap.String("page", page),
		zap.String("playlist", a.TargetPlaylist),
	)

	if err != nil {
		a.L.Error("failed to read playlist information to byte array",
			zap.String("page", page),
			zap.String("playlist", a.TargetPlaylist),
			zap.Error(err),
		)

		return nil
	}


	return toReturn
}