// Package downloader to handle downlader logic
package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aiyyra/song-manager/internal/playlist"
)

func Download(videoID string) error {
	stagingDir := "staging"
	url := "https://youtube.com/watch?v=" + videoID

	outputTemplate := filepath.Join(stagingDir, "%(title)s.%(ext)s")

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3", "-o", outputTemplate, url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Downloading %s...\n", videoID)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp download failed: %w", err)
	}

	fmt.Println("Download complete")
	return nil
}

func DownloadPlaylist(playlistID string) error {
	data, _ := playlist.Inspect(playlistID)

	for i, entry := range data.Entries {
		fmt.Printf("Downloading: %02d | %s | %s\n", i+1, entry.Title, entry.ID)
		if err := Download(entry.ID); err != nil {
			fmt.Printf("Download for `%s` failed\n", entry.Title)
		}
	}
	return nil
}
