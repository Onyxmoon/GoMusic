package model

import "time"

// Artist represents a music artist
type Artist struct {
	// Identification
	ID         string     `json:"id"`
	SourceID   string     `json:"sourceId"`
	SourceType SourceType `json:"sourceType"`

	// Metadata
	Name string `json:"name"`

	// Image
	ImagePath string `json:"imagePath,omitempty"`

	// Statistics
	AlbumCount int `json:"albumCount"`
	TrackCount int `json:"trackCount"`

	// Timestamps
	AddedAt time.Time `json:"addedAt"`
}