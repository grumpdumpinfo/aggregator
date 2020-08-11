package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func serverMode(app *App, port string) {
	// Will (hopefully) prevent us from getting hit with rate limits if somebody spams the server
	// Might have the unintended side effect of deadlocks though
	var limit int32
	// I don't think this will cause deadlocks, but if requests suddenly halt look for deadlocks here
	limitTicker := time.NewTicker(time.Second)
	done := make(chan bool)

	go func() {
		select {
		case <-done:
			return
		case <-limitTicker.C:
			atomic.StoreInt32(&limit, 0)
		}
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if int(atomic.LoadInt32(&limit)) < len(app.TargetPlaylists) {
			atomic.AddInt32(&limit, 1)
			go app.UpdateDatabase()
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
		}
	})

	server := http.Server{
		Addr:              port,
		Handler:           handler,
		ReadTimeout: 	   2 * time.Second,
		WriteTimeout:      2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	err := server.ListenAndServe()
	app.L.Error("server closed",
		zap.Error(err),
	)
}

func main() {
	// setting up logging
	l, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	// configuration
	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/gdiaggregator")

	viper.SetDefault("APIKey", "")
	viper.SetDefault("TargetPlaylists", "")
	viper.SetDefault("ServerPort", ":8080")

	err = viper.ReadInConfig()

	if err != nil {
		l.Panic("Failed to read configuration",
			zap.Error(err),
		)
	}

	apiKey := viper.GetString("APIKey")
	targetPlaylists := viper.GetStringSlice("TargetPlaylists")
	serverPort := viper.GetString("ServerPort")

	// TODO: validate inputs here to end execution early
	app := App{
		TargetPlaylists: targetPlaylists,
		APIKey:         apiKey,
		L:              l,
	}

	switch mode {
	case "build":
		app.BuildDatabase()
	case "update":
		app.UpdateDatabase()
	case "server":
		serverMode(&app, serverPort)
	default:
		fmt.Printf("unsupported mode %s, exiting...", mode)
		os.Exit(1)
	}
}
