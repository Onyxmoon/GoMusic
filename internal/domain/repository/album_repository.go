package repository

import (
	"context"

	"GoMusic/internal/domain/model"
)

// AlbumRepository defines the contract for album data access
type AlbumRepository interface {
	// CRUD operations
	FindByID(ctx context.Context, id string) (*model.Album, error)
	FindAll(ctx context.Context, opts *QueryOptions) ([]*model.Album, error)
	Create(ctx context.Context, album *model.Album) error
	Update(ctx context.Context, album *model.Album) error
	Delete(ctx context.Context, id string) error

	// Query operations
	FindByArtist(ctx context.Context, artistID string) ([]*model.Album, error)
	Search(ctx context.Context, query string, opts *SearchOptions) ([]*model.Album, error)

	// Source-specific operations
	GetSourceID() string
	GetSourceType() model.SourceType
}