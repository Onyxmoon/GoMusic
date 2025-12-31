package repository

import (
	"context"

	"GoMusic/internal/domain/model"
)

// PlaylistRepository defines the contract for playlist data access
type PlaylistRepository interface {
	// CRUD operations
	FindByID(ctx context.Context, id string) (*model.Playlist, error)
	FindAll(ctx context.Context) ([]*model.Playlist, error)
	Create(ctx context.Context, playlist *model.Playlist) error
	Update(ctx context.Context, playlist *model.Playlist) error
	Delete(ctx context.Context, id string) error

	// Playlist-specific operations
	AddTrack(ctx context.Context, playlistID, trackID string) error
	RemoveTrack(ctx context.Context, playlistID, trackID string) error
	ReorderTracks(ctx context.Context, playlistID string, trackIDs []string) error
	GetTracks(ctx context.Context, playlistID string) ([]*model.Track, error)
}