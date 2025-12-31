package filesystem

import "GoMusic/internal/domain/model"

// Extractor defines the interface for extracting metadata from audio files
type Extractor interface {
	// Extract reads metadata from an audio file and returns a Track
	Extract(filePath string) (*model.Track, error)

	// SupportsFormat checks if the extractor supports the given file extension
	SupportsFormat(extension string) bool
}

// Metadata holds extracted audio file metadata
type Metadata struct {
	// Basic info
	Title       string
	Artist      string
	Album       string
	AlbumArtist string
	Genre       string
	Year        int

	// Track info
	TrackNumber int
	DiscNumber  int

	// Audio properties
	Duration   int64 // in seconds
	BitRate    int   // kbps
	SampleRate int   // Hz
	Format     string

	// File info
	FilePath string
	FileSize int64

	// Artwork
	HasArtwork bool
	ArtworkExt string // jpg, png
}