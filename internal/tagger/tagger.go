// Package tagger tag the downloaded mp3 files
package tagger

import (
	"bytes"
	"fmt"

	"github.com/aiyyra/song-manager/internal/playlist"
	"github.com/bogem/id3v2/v2"
)

func ApplyTags(path string, track playlist.Track, playlistName string) error {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("failed to open mp3 for tagging: %w", err)
	}

	// Wipe Existing frames
	tag.DeleteAllFrames()

	// Basic Tags
	tag.SetTitle(track.Title)
	tag.SetArtist(joinArtist(track.Artists))
	tag.SetAlbum(playlistName)
	// tag.SetAlbumArtist(playlistName)
	// tag.SetTrack(fmt.Sprintf("%d", track.Position))

	commentFrame := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF16,
		Language:    "eng",
		Description: "source",
		Text:        "song-manager",
	}
	tag.AddCommentFrame(commentFrame)

	if err := tag.Save(); err != nil {
		return fmt.Errorf("failed to save tags: %w", err)
	}

	return tag.Close()
}

func joinArtist(artists []string) string {
	if len(artists) == 0 {
		return ""
	}
	var stringBuf bytes.Buffer
	stringBuf.WriteString(artists[0])
	for i := 1; i < len(artists); i++ {
		stringBuf.WriteString(", " + artists[i])
	}
	result := stringBuf.String()
	return result
}
