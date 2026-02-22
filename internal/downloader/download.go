// Package downloader to handle downlader logic
package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aiyyra/song-manager/internal/logger"
	"github.com/aiyyra/song-manager/internal/playlist"
	"github.com/aiyyra/song-manager/internal/tagger"
)

func Download(videoID string, track playlist.Track) error {
	stagingDir := "staging"
	url := "https://youtube.com/watch?v=" + videoID

	outputTemplate := filepath.Join(stagingDir, videoID+".%(ext)s")
	finalPath := filepath.Join(stagingDir, videoID+".mp3")

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3", "-o", outputTemplate, url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Downloading %s...\n", videoID)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp download failed: %w", err)
	}

	fmt.Println("Download complete")

	// replace playlist name with actual value (could extract or given as a flag).
	err := tagger.ApplyTags(finalPath, track, "temp")
	if err != nil {
		fmt.Printf("Taggig error: %s\n", err)
	}

	return nil
}

func DownloadPlaylist(playlistID string) error {
	// data, _ := playlist.Inspect(playlistID)
	fmt.Println("Starting download:")
	logger.Log.Info("Starting song-manager download-playlist ")
	tracks, _ := playlist.Inspect(playlistID)

	for _, entry := range tracks {
		fmt.Printf("Downloading: %02d | %s | %s\n", entry.Position, entry.Title, entry.VideoID)
		if err := Download(entry.VideoID, entry); err != nil {
			fmt.Printf("Download for `%s` failed\n", entry.Title)
		}
	}
	return nil
}
