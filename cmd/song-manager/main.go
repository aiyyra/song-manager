package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aiyyra/song-manager/internal/downloader"
	"github.com/aiyyra/song-manager/internal/playlist"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: song-manager <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "inspect":
		inspectCMD := flag.NewFlagSet("inspect", flag.ExitOnError)
		playlistID := inspectCMD.String("playlist", "", "Youtube playlist ID")

		if err := inspectCMD.Parse(os.Args[2:]); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if *playlistID == "" {
			fmt.Println("error: --playlist is required")
			os.Exit(1)
		}
		if err := playlist.Inspect(*playlistID); err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}

	case "download":
		{
			downloadCMD := flag.NewFlagSet("download", flag.ExitOnError)
			videoID := downloadCMD.String("video", "", "Video ID")

			if err := downloadCMD.Parse(os.Args[2:]); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}

			if *videoID == "" {
				fmt.Println("Error: --video is required")
				os.Exit(1)
			}
			if err := downloader.Download(*videoID); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		}

	default:
		fmt.Println("Unknown command: ", os.Args[1])
		os.Exit(1)
	}
}
