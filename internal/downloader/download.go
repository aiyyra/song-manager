// Package downloader to handle downlader logic
package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
