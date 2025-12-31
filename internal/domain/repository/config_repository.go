package repository

import (
	"context"

	"GoMusic/internal/domain/model"
)

// ConfigRepository defines the interface for configuration persistence
type ConfigRepository interface {
	// Load loads the application configuration
	Load(ctx context.Context) (*model.AppConfig, error)

	// Save saves the application configuration
	Save(ctx context.Context, config *model.AppConfig) error

	// Exists checks if a configuration file exists
	Exists(ctx context.Context) (bool, error)

	// Initialize creates a new default configuration
	Initialize(ctx context.Context) (*model.AppConfig, error)
}