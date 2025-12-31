package model

import "time"

// Track represents a music track from any source
type Track struct {
	// Identification
	ID         string     `json:"id"`
	SourceID   string     `json:"sourceId"`
	SourceType SourceType `json:"sourceType"`

	// Core metadata
	Title       string `json:"title"`
	Album       string `json:"album"`
	AlbumID     string `json:"albumId"`
	Artist      string `json:"artist"`
	ArtistID    string `json:"artistId"`
	AlbumArtist string `json:"albumArtist"`

	// Additional metadata
	Genre       string        `json:"genre"`
	Year        int           `json:"year"`
	TrackNumber int           `json:"trackNumber"`
	DiscNumber  int           `json:"discNumber"`
	Duration    time.Duration `json:"duration"`

	// File-specific (for filesystem sources)
	FilePath   string `json:"filePath,omitempty"`
	FileSize   int64  `json:"fileSize,omitempty"`
	Format     string `json:"format"`           // mp3, flac, m4a, ogg, etc.
	BitRate    int    `json:"bitRate,omitempty"` // kbps
	SampleRate int    `json:"sampleRate,omitempty"` // Hz

	// API-specific (for API sources)
	ExternalID string `json:"externalId,omitempty"`
	StreamURL  string `json:"streamUrl,omitempty"`

	// Artwork
	ArtworkPath string `json:"artworkPath,omitempty"` // Local cache path or URL

	// Timestamps
	AddedAt    time.Time `json:"addedAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}