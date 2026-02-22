package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aiyyra/song-manager/internal/downloader"
	"github.com/aiyyra/song-manager/internal/logger"
	"github.com/aiyyra/song-manager/internal/playlist"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Application panicked", zap.Any("panic", r))
		}

		if err := logger.Log.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
		}
	}()

	logger.Log.Info("Starting song-manager")

	if len(os.Args) < 2 {
		fmt.Println("Usage: song-manager <command> [options]")
		os.Exit(1)
	}

	os.Exit(run())
}

func run() int {
	logger.Log.Info("Starting song-manager")

	if len(os.Args) < 2 {
		logger.Log.Error("Missing command")
		return 1
	}

	switch os.Args[1] {
	case "inspect":
		return handleInspect()
	case "fetch":
		return handleFetch()
	default:
		logger.Log.Error("Unknown command",
			zap.String("command", os.Args[1]),
		)
		return 1
	}
}

func handleFetch() int {
	fetchCMD := flag.NewFlagSet("fetch", flag.ExitOnError)
	playlistID := fetchCMD.String("playlist", "", "Youtube Playlist ID")

	if err := fetchCMD.Parse(os.Args[2:]); err != nil {
		logger.Log.Error("Error parsing fetch")
		return 1
	}

	if *playlistID == "" {
		logger.Log.Error("error: --video is required")
		return 1
	}
	if err := downloader.DownloadPlaylist(*playlistID); err != nil {
		logger.Log.Error("Error calling DownloadPlaylist")
		return 1
	}
	return 0
}

func handleInspect() int {
	inspectCMD := flag.NewFlagSet("inspect", flag.ExitOnError)
	playlistID := inspectCMD.String("playlist", "", "Youtube playlist ID")

	if err := inspectCMD.Parse(os.Args[2:]); err != nil {
		logger.Log.Error("Error parsing inspect")
		return 1
	}

	if *playlistID == "" {
		logger.Log.Error("error: --playlist is required")
		return 1
	}
	if _, err := playlist.Inspect(*playlistID); err != nil {
		logger.Log.Error("Error calling inspect")
		return 1
	}
	return 0
}
