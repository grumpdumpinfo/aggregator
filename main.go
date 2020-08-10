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
	viper.SetDefault("TargetPlaylist", "")

	err = viper.ReadInConfig()

	if err != nil {
		l.Panic("Failed to read configuration",
			zap.Error(err),
		)
	}

	apiKey := viper.GetString("APIKey")
	targetPlaylist := viper.GetString("TargetPlaylist")

	// TODO: validate inputs here to end execution early
	app := App{
		TargetPlaylist: targetPlaylist,
		APIKey:         apiKey,
		L:              l,
	}

	firstPage := string(app.GetPlaylistInfoPage(""))

	fmt.Printf(firstPage)
}
