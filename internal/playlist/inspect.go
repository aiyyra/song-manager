// Package playlist provides functionality for inspecting anf managing playlist
package playlist

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/aiyyra/song-manager/internal/logger"
	"go.uber.org/zap"
)

type YTDLPPlaylist struct {
	Title   string `json:"title"`
	ID      string `json:"id"`
	Entries []struct {
		Entry
	} `json:"entries"`
}

type Entry struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Uploader    string  `json:"uploader"`
	Channel     string  `json:"channel"`
	Artist      string  `json:"artist"`
	Creator     string  `json:"creator"`
	Duration    int     `json:"duration"`
	UploadDate  string  `json:"upload_date"`
	Description string  `json:"description"`
	Track       string  `json:"track"`
	Album       string  `json:"album"`
	Thumbnail   string  `json:"thumbnail"`
	Thumbnails  []Thumb `json:"thumbnails"`
}

type Thumb struct {
	URL string `json:"url"`
}

type Track struct {
	VideoID  string
	Title    string
	Artists  []string
	Duration int
	Position int
}

func Inspect(playlistID string) ([]Track, error) {
	logger.Log.Info("Inspecting playlist",
		zap.String("playlistID", playlistID),
	)
	url := "https://music.youtube.com/playlist?list=" + playlistID
	// url := "https://www.youtube.com/playlist/list=" + playlistID

	logger.Log.Debug("Executing yt-dlp", zap.String("url", url))

	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", url)

	output, err := cmd.Output()
	if err != nil {
		logger.Log.Error("yt-dlp execution failed", zap.String("playlistID", playlistID), zap.Error(err))
		return nil, fmt.Errorf("failed to execute yt-dlp: %w", err)
	}
	logger.Log.Debug("yt-dlp output received",
		zap.Int("bytes", len(output)),
	)

	var data YTDLPPlaylist
	if err := json.Unmarshal(output, &data); err != nil {
		logger.Log.Error("Failed to parse yt-dlp JSON",
			zap.String("playlist_id", playlistID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed ro parse yt-dlp JSON: %w", err)
	}
	logger.Log.Info("Playlist parsed successfully",
		zap.Int("track_count", len(data.Entries)),
	)

	for i, entry := range data.Entries {
		fmt.Printf("%02d | %s | %s\n", i+1, entry.Title, entry.ID)
	}

	tracks := ConvertToTrack(data)
	logger.Log.Debug("Tracks converted",
		zap.Int("track_count", len(tracks)),
	)

	return tracks, nil
}

func ConvertToTrack(p YTDLPPlaylist) []Track {
	tracks := make([]Track, 0, len(p.Entries))

	for i, entry := range p.Entries {
		tracks = append(tracks, Track{
			VideoID:  entry.ID,
			Title:    entry.Title,
			Artists:  resolveArtist(entry),
			Duration: entry.Duration,
			Position: i + 1,
		})
	}
	return tracks
}

func resolveArtist(e struct{ Entry }) []string {
	if e.Artist != "" {
		return []string{e.Artist}
	}
	if e.Creator != "" {
		return []string{e.Creator}
	}
	if e.Uploader != "" {
		return []string{e.Uploader}
	}
	if e.Channel != "" {
		return []string{e.Channel}
	}
	return []string{"Unknown Artist"}
}
