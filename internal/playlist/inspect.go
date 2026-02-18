// Package playlist provides functionality for inspecting anf managing playlist
package playlist

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type YTDLPPlaylist struct {
	Entries []struct {
		Entry
	} `json:"entries"`
}

type Entry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func Inspect(playlistID string) error {
	url := "https://music.youtube.com/playlist?list=" + playlistID

	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", url)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute yt-dlp: %w", err)
	}

	var data YTDLPPlaylist
	if err := json.Unmarshal(output, &data); err != nil {
		return fmt.Errorf("failed ro parse yt-dlp JSON: %w", err)
	}

	for i, entry := range data.Entries {
		fmt.Printf("%02d | %s | %s\n", i+1, entry.Title, entry.ID)
	}

	return nil
}
