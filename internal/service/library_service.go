package service

import (
	"context"
	"sync"

	"GoMusic/internal/domain/model"
	"GoMusic/internal/domain/repository"
	"GoMusic/internal/util/errors"
)

// LibraryService aggregates multiple track repositories (sources)
type LibraryService struct {
	trackRepos    map[string]repository.TrackRepository
	albumRepos    map[string]repository.AlbumRepository
	artistRepos   map[string]repository.ArtistRepository
	playlistRepo  repository.PlaylistRepository
	mu            sync.RWMutex
}

// NewLibraryService creates a new library service
func NewLibraryService() *LibraryService {
	return &LibraryService{
		trackRepos:  make(map[string]repository.TrackRepository),
		albumRepos:  make(map[string]repository.AlbumRepository),
		artistRepos: make(map[string]repository.ArtistRepository),
	}
}

// RegisterTrackRepository adds a track repository to the library
func (s *LibraryService) RegisterTrackRepository(sourceID string, repo repository.TrackRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trackRepos[sourceID] = repo
}

// UnregisterTrackRepository removes a track repository from the library
func (s *LibraryService) UnregisterTrackRepository(sourceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.trackRepos, sourceID)
}

// RegisterAlbumRepository adds an album repository to the library
func (s *LibraryService) RegisterAlbumRepository(sourceID string, repo repository.AlbumRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.albumRepos[sourceID] = repo
}

// RegisterArtistRepository adds an artist repository to the library
func (s *LibraryService) RegisterArtistRepository(sourceID string, repo repository.ArtistRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.artistRepos[sourceID] = repo
}

// SetPlaylistRepository sets the playlist repository
func (s *LibraryService) SetPlaylistRepository(repo repository.PlaylistRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.playlistRepo = repo
}

// GetAllTracks retrieves tracks from all sources
func (s *LibraryService) GetAllTracks(ctx context.Context, opts *repository.QueryOptions) ([]*model.Track, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allTracks []*model.Track

	for _, repo := range s.trackRepos {
		tracks, err := repo.FindAll(ctx, opts)
		if err != nil {
			// Log error but continue with other sources
			continue
		}
		allTracks = append(allTracks, tracks...)
	}

	return allTracks, nil
}

// GetTrackByID searches all sources for a track
func (s *LibraryService) GetTrackByID(ctx context.Context, id string) (*model.Track, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, repo := range s.trackRepos {
		track, err := repo.FindByID(ctx, id)
		if err == nil {
			return track, nil
		}
	}

	return nil, errors.ErrNotFound
}

// SearchTracks searches for tracks across all sources
func (s *LibraryService) SearchTracks(ctx context.Context, query string, opts *repository.SearchOptions) ([]*model.Track, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allResults []*model.Track

	for _, repo := range s.trackRepos {
		results, err := repo.Search(ctx, query, opts)
		if err != nil {
			// Log error but continue with other sources
			continue
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// GetTracksByAlbum retrieves tracks for a specific album
func (s *LibraryService) GetTracksByAlbum(ctx context.Context, albumID string) ([]*model.Track, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allTracks []*model.Track

	for _, repo := range s.trackRepos {
		tracks, err := repo.FindByAlbum(ctx, albumID)
		if err != nil {
			continue
		}
		allTracks = append(allTracks, tracks...)
	}

	if len(allTracks) == 0 {
		return nil, errors.ErrNotFound
	}

	return allTracks, nil
}

// GetTracksByArtist retrieves tracks for a specific artist
func (s *LibraryService) GetTracksByArtist(ctx context.Context, artistID string) ([]*model.Track, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allTracks []*model.Track

	for _, repo := range s.trackRepos {
		tracks, err := repo.FindByArtist(ctx, artistID)
		if err != nil {
			continue
		}
		allTracks = append(allTracks, tracks...)
	}

	if len(allTracks) == 0 {
		return nil, errors.ErrNotFound
	}

	return allTracks, nil
}

// ScanSource triggers a scan on a specific source
func (s *LibraryService) ScanSource(ctx context.Context, sourceID string) error {
	s.mu.RLock()
	repo, exists := s.trackRepos[sourceID]
	s.mu.RUnlock()

	if !exists {
		return errors.ErrSourceNotFound
	}

	return repo.Scan(ctx)
}

// ScanAllSources triggers a scan on all sources
func (s *LibraryService) ScanAllSources(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var wg sync.WaitGroup
	errCh := make(chan error, len(s.trackRepos))

	for sourceID, repo := range s.trackRepos {
		wg.Add(1)
		go func(id string, r repository.TrackRepository) {
			defer wg.Done()
			if err := r.Scan(ctx); err != nil {
				errCh <- err
			}
		}(sourceID, repo)
	}

	wg.Wait()
	close(errCh)

	// Collect any errors (returns first error encountered)
	for err := range errCh {
		return err
	}

	return nil
}

// GetScanProgress retrieves scan progress for a specific source
func (s *LibraryService) GetScanProgress(sourceID string) (*repository.ScanProgress, error) {
	s.mu.RLock()
	repo, exists := s.trackRepos[sourceID]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.ErrSourceNotFound
	}

	return repo.GetScanProgress(), nil
}

// GetAllScanProgress retrieves scan progress for all sources
func (s *LibraryService) GetAllScanProgress() map[string]*repository.ScanProgress {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress := make(map[string]*repository.ScanProgress)

	for sourceID, repo := range s.trackRepos {
		progress[sourceID] = repo.GetScanProgress()
	}

	return progress
}

// GetSources returns information about all registered sources
func (s *LibraryService) GetSources() []SourceInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sources := make([]SourceInfo, 0, len(s.trackRepos))

	for sourceID, repo := range s.trackRepos {
		sources = append(sources, SourceInfo{
			ID:   sourceID,
			Type: repo.GetSourceType(),
		})
	}

	return sources
}

// SourceInfo contains basic information about a source
type SourceInfo struct {
	ID   string           `json:"id"`
	Type model.SourceType `json:"type"`
}

// GetRepositories returns all registered track repositories
func (s *LibraryService) GetRepositories() map[string]repository.TrackRepository {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	repos := make(map[string]repository.TrackRepository)
	for id, repo := range s.trackRepos {
		repos[id] = repo
	}

	return repos
}
