package filesystem

import (
	"sort"
	"strings"
	"sync"

	"GoMusic/internal/domain/model"
	"GoMusic/internal/domain/repository"
)

// TrackCache provides thread-safe in-memory storage for tracks
type TrackCache struct {
	tracks map[string]*model.Track
	mu     sync.RWMutex
}

// NewTrackCache creates a new track cache
func NewTrackCache() *TrackCache {
	return &TrackCache{
		tracks: make(map[string]*model.Track),
	}
}

// Add adds a track to the cache
func (c *TrackCache) Add(track *model.Track) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tracks[track.ID] = track
}

// Get retrieves a track by ID
func (c *TrackCache) Get(id string) *model.Track {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tracks[id]
}

// GetAll retrieves all tracks with optional filtering and pagination
func (c *TrackCache) GetAll(opts *repository.QueryOptions) []*model.Track {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if opts == nil {
		opts = repository.DefaultQueryOptions()
	}

	// Collect all tracks
	allTracks := make([]*model.Track, 0, len(c.tracks))
	for _, track := range c.tracks {
		allTracks = append(allTracks, track)
	}

	// Apply sorting
	c.sortTracks(allTracks, opts.SortBy, opts.SortOrder)

	// Apply pagination
	start := opts.Offset
	if start > len(allTracks) {
		return []*model.Track{}
	}

	// Calculate end position (Limit=0 means "all tracks")
	end := start + opts.Limit
	if opts.Limit == 0 {
		end = len(allTracks)
	}
	if end > len(allTracks) {
		end = len(allTracks)
	}

	return allTracks[start:end]
}

// Search searches for tracks matching a query string
func (c *TrackCache) Search(query string, opts *repository.SearchOptions) []*model.Track {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if opts == nil {
		opts = repository.DefaultSearchOptions()
	}

	query = strings.ToLower(query)
	var results []*model.Track

	for _, track := range c.tracks {
		if c.matchesQuery(track, query, opts.Fields) {
			results = append(results, track)
		}
	}

	// Apply sorting
	c.sortTracks(results, opts.SortBy, opts.SortOrder)

	// Apply pagination
	start := opts.Offset
	if start > len(results) {
		return []*model.Track{}
	}

	end := start + opts.Limit
	if opts.Limit == 0 {
		end = len(results)
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end]
}

// FindByAlbum returns all tracks for a given album ID
func (c *TrackCache) FindByAlbum(albumID string) []*model.Track {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var tracks []*model.Track
	for _, track := range c.tracks {
		if track.AlbumID == albumID {
			tracks = append(tracks, track)
		}
	}

	// Sort by disc and track number
	sort.Slice(tracks, func(i, j int) bool {
		if tracks[i].DiscNumber != tracks[j].DiscNumber {
			return tracks[i].DiscNumber < tracks[j].DiscNumber
		}
		return tracks[i].TrackNumber < tracks[j].TrackNumber
	})

	return tracks
}

// FindByArtist returns all tracks for a given artist ID
func (c *TrackCache) FindByArtist(artistID string) []*model.Track {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var tracks []*model.Track
	for _, track := range c.tracks {
		if track.ArtistID == artistID {
			tracks = append(tracks, track)
		}
	}

	// Sort by album and track number
	c.sortTracks(tracks, "album", "asc")

	return tracks
}

// Delete removes a track from the cache
func (c *TrackCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.tracks, id)
}

// Clear removes all tracks from the cache
func (c *TrackCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tracks = make(map[string]*model.Track)
}

// Count returns the total number of tracks in the cache
func (c *TrackCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.tracks)
}

// matchesQuery checks if a track matches the search query in any of the specified fields
func (c *TrackCache) matchesQuery(track *model.Track, query string, fields []string) bool {
	if len(fields) == 0 {
		fields = []string{"title", "artist", "album"}
	}

	for _, field := range fields {
		var value string
		switch field {
		case "title":
			value = track.Title
		case "artist":
			value = track.Artist
		case "album":
			value = track.Album
		case "genre":
			value = track.Genre
		}

		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}

	return false
}

// sortTracks sorts a slice of tracks based on the given field and order
func (c *TrackCache) sortTracks(tracks []*model.Track, sortBy, sortOrder string) {
	if sortBy == "" {
		sortBy = "title"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	sort.Slice(tracks, func(i, j int) bool {
		var less bool

		switch sortBy {
		case "title":
			less = tracks[i].Title < tracks[j].Title
		case "artist":
			less = tracks[i].Artist < tracks[j].Artist
		case "album":
			less = tracks[i].Album < tracks[j].Album
		case "year":
			less = tracks[i].Year < tracks[j].Year
		case "duration":
			less = tracks[i].Duration < tracks[j].Duration
		case "addedAt":
			less = tracks[i].AddedAt.Before(tracks[j].AddedAt)
		default:
			less = tracks[i].Title < tracks[j].Title
		}

		if sortOrder == "desc" {
			return !less
		}
		return less
	})
}