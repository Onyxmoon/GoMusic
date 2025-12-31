package model

import "time"

// Album represents a music album
type Album struct {
	// Identification
	ID         string     `json:"id"`
	SourceID   string     `json:"sourceId"`
	SourceType SourceType `json:"sourceType"`

	// Metadata
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	ArtistID  string `json:"artistId"`
	Year      int    `json:"year"`
	Genre     string `json:"genre"`

	// Artwork
	ArtworkPath string `json:"artworkPath,omitempty"`

	// Statistics
	TrackCount    int           `json:"trackCount"`
	TotalDuration time.Duration `json:"totalDuration"`

	// Timestamps
	AddedAt time.Time `json:"addedAt"`
}