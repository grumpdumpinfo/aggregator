package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// setting up logging
	l, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	// configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/gdiaggregator")

	viper.SetDefault("APIKey", "")
	viper.SetDefault("TargetPlaylists", "")

	err = viper.ReadInConfig()

	if err != nil {
		l.Panic("Failed to read configuration",
			zap.Error(err),
		)
	}

	apiKey := viper.GetString("APIKey")
	targetPlaylists := viper.GetStringSlice("TargetPlaylists")

	// TODO: validate inputs here to end execution early
	app := App{
		TargetPlaylists: targetPlaylists,
		APIKey:         apiKey,
		L:              l,
	}

	l.Info("getting all playlists first pages",
		zap.Strings("targetPlaylists", app.TargetPlaylists)	,
	)

	pages := app.GetAllPlaylistsPages()

	for _, page := range pages {
		fmt.Printf("%d\n", len(page.Items))
	}
}
