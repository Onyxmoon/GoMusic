package filesystem

import (
	"context"
	"fmt"
	"sync"

	"GoMusic/internal/domain/model"
	"GoMusic/internal/domain/repository"
	"GoMusic/internal/domain/source/capability"
	"GoMusic/internal/util/errors"
)

// filesystemTrackRepository implements TrackRepository for local filesystem
// This struct is unexported to enforce usage of the constructor
type filesystemTrackRepository struct {
	sourceID     string
	config       *model.FilesystemSourceConfig
	cache        *TrackCache
	scanner      *DirectoryScanner
	extractor    Extractor
	scanProgress *repository.ScanProgress
	mu           sync.RWMutex
}

// NewFilesystemTrackRepository creates a new filesystem track repository
func NewFilesystemTrackRepository(
	sourceID string,
	config *model.FilesystemSourceConfig,
	extractor Extractor,
) repository.TrackRepository {
	return &filesystemTrackRepository{
		sourceID:  sourceID,
		config:    config,
		cache:     NewTrackCache(),
		scanner:   NewDirectoryScanner(config.RootPath, config.SupportedFormats),
		extractor: extractor,
		scanProgress: &repository.ScanProgress{
			IsScanning: false,
		},
	}
}

// FindByID finds a track by ID
func (r *filesystemTrackRepository) FindByID(ctx context.Context, id string) (*model.Track, error) {
	track := r.cache.Get(id)
	if track == nil {
		return nil, errors.ErrNotFound
	}
	return track, nil
}

// FindAll returns all tracks with optional filtering
func (r *filesystemTrackRepository) FindAll(ctx context.Context, opts *repository.QueryOptions) ([]*model.Track, error) {
	return r.cache.GetAll(opts), nil
}

// Create adds a new track to the repository
func (r *filesystemTrackRepository) Create(ctx context.Context, track *model.Track) error {
	if r.cache.Get(track.ID) != nil {
		return errors.ErrAlreadyExists
	}
	r.cache.Add(track)
	return nil
}

// Update updates an existing track
func (r *filesystemTrackRepository) Update(ctx context.Context, track *model.Track) error {
	if r.cache.Get(track.ID) == nil {
		return errors.ErrNotFound
	}
	r.cache.Add(track)
	return nil
}

// Delete removes a track from the repository
func (r *filesystemTrackRepository) Delete(ctx context.Context, id string) error {
	if r.cache.Get(id) == nil {
		return errors.ErrNotFound
	}
	r.cache.Delete(id)
	return nil
}

// FindByAlbum returns all tracks for a given album
func (r *filesystemTrackRepository) FindByAlbum(ctx context.Context, albumID string) ([]*model.Track, error) {
	return r.cache.FindByAlbum(albumID), nil
}

// FindByArtist returns all tracks for a given artist
func (r *filesystemTrackRepository) FindByArtist(ctx context.Context, artistID string) ([]*model.Track, error) {
	return r.cache.FindByArtist(artistID), nil
}

// Search searches for tracks matching the query
func (r *filesystemTrackRepository) Search(ctx context.Context, query string, opts *repository.SearchOptions) ([]*model.Track, error) {
	return r.cache.Search(query, opts), nil
}

// GetSourceID returns the source ID
func (r *filesystemTrackRepository) GetSourceID() string {
	return r.sourceID
}

// GetSourceType returns the source type
func (r *filesystemTrackRepository) GetSourceType() model.SourceType {
	return model.SourceTypeFilesystem
}

// Scan scans the filesystem for audio files and extracts metadata
func (r *filesystemTrackRepository) Scan(ctx context.Context) error {
	r.mu.Lock()
	if r.scanProgress.IsScanning {
		r.mu.Unlock()
		return errors.ErrScanInProgress
	}
	r.scanProgress.IsScanning = true
	r.scanProgress.ProcessedFiles = 0
	r.scanProgress.TotalFiles = 0
	r.scanProgress.CurrentFile = ""
	r.scanProgress.Errors = []string{}
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.scanProgress.IsScanning = false
		r.scanProgress.CurrentFile = ""
		r.mu.Unlock()
	}()

	// Clear existing cache
	r.cache.Clear()

	// Scan directory for audio files
	files, err := r.scanner.ScanDirectory(ctx, func(filePath string) {
		r.mu.Lock()
		r.scanProgress.CurrentFile = filePath
		r.scanProgress.ProcessedFiles++
		r.mu.Unlock()
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	r.mu.Lock()
	r.scanProgress.TotalFiles = len(files)
	r.scanProgress.ProcessedFiles = 0
	r.mu.Unlock()

	// Extract metadata from each file
	for _, filePath := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			r.mu.Lock()
			r.scanProgress.CurrentFile = filePath
			r.scanProgress.ProcessedFiles++
			r.mu.Unlock()

			track, err := r.extractor.Extract(filePath)
			if err != nil {
				errMsg := fmt.Sprintf("%s: %v", filePath, err)
				r.mu.Lock()
				r.scanProgress.Errors = append(r.scanProgress.Errors, errMsg)
				r.mu.Unlock()
				continue
			}

			// Ensure track has the correct source ID
			track.SourceID = r.sourceID
			track.SourceType = model.SourceTypeFilesystem

			r.cache.Add(track)
		}
	}

	return nil
}

// GetScanProgress returns the current scan progress
func (r *filesystemTrackRepository) GetScanProgress() *repository.ScanProgress {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to avoid data races
	progress := *r.scanProgress
	progress.Errors = make([]string, len(r.scanProgress.Errors))
	copy(progress.Errors, r.scanProgress.Errors)

	return &progress
}

// GetRootPath returns the configured root path for this repository
func (r *filesystemTrackRepository) GetRootPath() string {
	return r.config.RootPath
}

// ListDirectory lists the contents of a directory
// relativePath is relative to the root path ("" or "/" for root)
func (r *filesystemTrackRepository) ListDirectory(relativePath string) ([]*capability.FileNode, error) {
	// Build full path
	fullPath := r.config.RootPath
	if relativePath != "" && relativePath != "/" {
		fullPath = fmt.Sprintf("%s%s", r.config.RootPath, relativePath)
	}

	// Read directory contents
	entries, err := r.scanner.ListDirectory(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	return entries, nil
}
