package repository

import (
	"context"

	"GoMusic/internal/domain/model"
)

// TrackRepository defines the contract for track data access
type TrackRepository interface {
	// CRUD operations
	FindByID(ctx context.Context, id string) (*model.Track, error)
	FindAll(ctx context.Context, opts *QueryOptions) ([]*model.Track, error)
	Create(ctx context.Context, track *model.Track) error
	Update(ctx context.Context, track *model.Track) error
	Delete(ctx context.Context, id string) error

	// Query operations
	FindByAlbum(ctx context.Context, albumID string) ([]*model.Track, error)
	FindByArtist(ctx context.Context, artistID string) ([]*model.Track, error)
	Search(ctx context.Context, query string, opts *SearchOptions) ([]*model.Track, error)

	// Source-specific operations
	GetSourceID() string
	GetSourceType() model.SourceType

	// Scanning (for repositories that support it, like filesystem)
	Scan(ctx context.Context) error
	GetScanProgress() *ScanProgress
}

// QueryOptions provides pagination and filtering for queries
type QueryOptions struct {
	Limit     int                    `json:"limit"`
	Offset    int                    `json:"offset"`
	SortBy    string                 `json:"sortBy"`    // Field to sort by
	SortOrder string                 `json:"sortOrder"` // "asc" or "desc"
	Filters   map[string]interface{} `json:"filters"`   // Generic filters
}

// SearchOptions extends QueryOptions with search-specific options
type SearchOptions struct {
	*QueryOptions
	Fields []string `json:"fields"` // Which fields to search (title, artist, album, etc.)
}

// ScanProgress tracks the progress of a repository scan operation
type ScanProgress struct {
	IsScanning     bool     `json:"isScanning"`
	TotalFiles     int      `json:"totalFiles"`
	ProcessedFiles int      `json:"processedFiles"`
	CurrentFile    string   `json:"currentFile"`
	Errors         []string `json:"errors,omitempty"`
}

// DefaultQueryOptions returns default query options
func DefaultQueryOptions() *QueryOptions {
	return &QueryOptions{
		Limit:     0,
		Offset:    0,
		SortBy:    "title",
		SortOrder: "asc",
		Filters:   make(map[string]interface{}),
	}
}

// DefaultSearchOptions returns default search options
func DefaultSearchOptions() *SearchOptions {
	return &SearchOptions{
		QueryOptions: DefaultQueryOptions(),
		Fields:       []string{"title", "artist", "album"},
	}
}
