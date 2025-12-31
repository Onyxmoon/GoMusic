package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"

	"GoMusic/internal/domain/model"
)

// TagExtractor uses github.com/dhowden/tag to extract metadata
// Supports ID3 (MP3, M4A) and Vorbis (FLAC, OGG) formats
type TagExtractor struct {
	sourceID string
}

// NewTagExtractor creates a new tag-based metadata extractor
func NewTagExtractor(sourceID string) *TagExtractor {
	return &TagExtractor{
		sourceID: sourceID,
	}
}

// SupportsFormat checks if the file format is supported
func (e *TagExtractor) SupportsFormat(extension string) bool {
	ext := strings.ToLower(extension)
	switch ext {
	case ".mp3", ".m4a", ".m4b", ".m4p", ".flac", ".ogg", ".oga":
		return true
	default:
		return false
	}
}

// Extract reads metadata from an audio file
func (e *TagExtractor) Extract(filePath string) (*model.Track, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Read metadata tags
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		// If tag reading fails, create a track with basic info from filename
		return e.createTrackFromFilename(filePath, fileInfo), nil
	}

	// Create track from metadata
	track := &model.Track{
		ID:         generateTrackID(filePath),
		SourceID:   e.sourceID,
		SourceType: model.SourceTypeFilesystem,

		// Metadata from tags
		Title:       metadata.Title(),
		Artist:      metadata.Artist(),
		Album:       metadata.Album(),
		AlbumArtist: metadata.AlbumArtist(),
		Genre:       metadata.Genre(),
		Year:        metadata.Year(),

		// Track numbers
		TrackNumber: getTrackNumber(metadata),
		DiscNumber:  getDiscNumber(metadata),

		// File information
		FilePath: filePath,
		FileSize: fileInfo.Size(),
		Format:   getFormat(metadata),

		// Audio properties
		BitRate:    getBitRate(metadata),
		SampleRate: getSampleRate(metadata),

		// Timestamps
		AddedAt:    time.Now(),
		ModifiedAt: fileInfo.ModTime(),
	}

	// Extract and save artwork if available
	if picture := metadata.Picture(); picture != nil {
		artworkPath, err := e.saveArtwork(track.ID, picture)
		if err != nil {
			// Log error but don't fail the entire extraction
			fmt.Printf("Failed to save artwork for %s: %v\n", track.Title, err)
		} else {
			track.ArtworkPath = artworkPath
		}
	}

	// Extract audio properties (duration, sample rate, bitrate) from file
	// This is specific to filesystem sources - API sources get this from the API
	analyzer := NewAudioAnalyzer()
	props := analyzer.Analyze(filePath)
	if props.Duration > 0 {
		track.Duration = props.Duration
	}
	if props.SampleRate > 0 {
		track.SampleRate = props.SampleRate
	}
	if props.BitRate > 0 {
		track.BitRate = props.BitRate
	}

	// Use filename as title if title is empty
	if track.Title == "" {
		track.Title = getFilenameWithoutExt(filePath)
	}

	// Use "Unknown Artist" if artist is empty
	if track.Artist == "" {
		track.Artist = "Unknown Artist"
	}

	// Use "Unknown Album" if album is empty
	if track.Album == "" {
		track.Album = "Unknown Album"
	}

	// Generate IDs for album and artist
	track.AlbumID = generateAlbumID(track.Album, track.AlbumArtist)
	track.ArtistID = generateArtistID(track.Artist)

	return track, nil
}

// createTrackFromFilename creates a basic track when metadata extraction fails
func (e *TagExtractor) createTrackFromFilename(filePath string, fileInfo os.FileInfo) *model.Track {
	return &model.Track{
		ID:         generateTrackID(filePath),
		SourceID:   e.sourceID,
		SourceType: model.SourceTypeFilesystem,
		Title:      getFilenameWithoutExt(filePath),
		Artist:     "Unknown Artist",
		Album:      "Unknown Album",
		FilePath:   filePath,
		FileSize:   fileInfo.Size(),
		Format:     strings.TrimPrefix(filepath.Ext(filePath), "."),
		ArtistID:   generateArtistID("Unknown Artist"),
		AlbumID:    generateAlbumID("Unknown Album", ""),
		AddedAt:    time.Now(),
		ModifiedAt: fileInfo.ModTime(),
	}
}

// Helper functions

func getTrackNumber(m tag.Metadata) int {
	track, _ := m.Track()
	return track
}

func getDiscNumber(m tag.Metadata) int {
	disc, _ := m.Disc()
	return disc
}

func getFormat(m tag.Metadata) string {
	return strings.ToLower(string(m.Format()))
}

// TODO: Implement bitrate extraction
func getBitRate(m tag.Metadata) int {
	return 0
}

// TODO: Implement sample rate extraction
func getSampleRate(m tag.Metadata) int {
	return 0
}

// TODO: Implement duration extraction
// github.com/dhowden/tag provides duration through FileType interface
// but it's not always available - needs proper implementation
func getDuration(m tag.Metadata) int64 {
	return 0
}

func getFilenameWithoutExt(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// ID generation functions

func generateTrackID(filePath string) string {
	hash := sha256.Sum256([]byte(filePath))
	return "track_" + hex.EncodeToString(hash[:8])
}

func generateAlbumID(album, albumArtist string) string {
	key := album
	if albumArtist != "" {
		key = albumArtist + "_" + album
	}
	hash := sha256.Sum256([]byte(key))
	return "album_" + hex.EncodeToString(hash[:8])
}

func generateArtistID(artist string) string {
	hash := sha256.Sum256([]byte(artist))
	return "artist_" + hex.EncodeToString(hash[:8])
}

// saveArtwork saves album artwork to the cache directory
func (e *TagExtractor) saveArtwork(trackID string, picture *tag.Picture) (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create artwork cache directory: ~/.gomusic/artwork/
	artworkDir := filepath.Join(homeDir, ".gomusic", "artwork")
	if err := os.MkdirAll(artworkDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create artwork directory: %w", err)
	}

	// Determine file extension from MIME type
	ext := getImageExtension(picture.MIMEType)

	// Generate filename: <trackID>.<ext>
	filename := trackID + ext
	artworkPath := filepath.Join(artworkDir, filename)

	// Write artwork data to file
	if err := os.WriteFile(artworkPath, picture.Data, 0644); err != nil {
		return "", fmt.Errorf("failed to write artwork file: %w", err)
	}

	return filename, nil
}

// getImageExtension returns the file extension for a given MIME type
func getImageExtension(mimeType string) string {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/bmp":
		return ".bmp"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg" // Default to JPEG
	}
}