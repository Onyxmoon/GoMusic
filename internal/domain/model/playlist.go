package model

import "time"

// Playlist represents a user-created playlist
type Playlist struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TrackIDs    []string `json:"trackIds"` // Ordered list of track IDs

	// Artwork (generated from tracks or custom)
	CoverPath string `json:"coverPath,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}