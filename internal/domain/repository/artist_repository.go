package repository

import (
	"context"

	"GoMusic/internal/domain/model"
)

// ArtistRepository defines the contract for artist data access
type ArtistRepository interface {
	// CRUD operations
	FindByID(ctx context.Context, id string) (*model.Artist, error)
	FindAll(ctx context.Context, opts *QueryOptions) ([]*model.Artist, error)
	Create(ctx context.Context, artist *model.Artist) error
	Update(ctx context.Context, artist *model.Artist) error
	Delete(ctx context.Context, id string) error

	// Query operations
	Search(ctx context.Context, query string, opts *SearchOptions) ([]*model.Artist, error)

	// Source-specific operations
	GetSourceID() string
	GetSourceType() model.SourceType
}