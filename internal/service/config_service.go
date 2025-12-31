package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"GoMusic/internal/domain/model"
	"GoMusic/internal/domain/repository"
	"GoMusic/internal/util/errors"
)

// ConfigService manages application configuration
type ConfigService struct {
	repo   repository.ConfigRepository
	config *model.AppConfig
	mu     sync.RWMutex
}

// NewConfigService creates a new config service
func NewConfigService(repo repository.ConfigRepository) *ConfigService {
	return &ConfigService{
		repo: repo,
	}
}

// Initialize loads or creates the configuration
func (s *ConfigService) Initialize(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if config exists
	exists, err := s.repo.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check config existence: %w", err)
	}

	// Load existing or create new
	if exists {
		config, err := s.repo.Load(ctx)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		s.config = config
	} else {
		config, err := s.repo.Initialize(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		s.config = config
	}

	return nil
}

// GetConfig returns a copy of the current configuration
func (s *ConfigService) GetConfig() *model.AppConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modifications
	configCopy := *s.config
	configCopy.Sources = make([]model.SourceConfiguration, len(s.config.Sources))
	copy(configCopy.Sources, s.config.Sources)

	return &configCopy
}

// GetSources returns all configured sources
func (s *ConfigService) GetSources() []model.SourceConfiguration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sources := make([]model.SourceConfiguration, len(s.config.Sources))
	copy(sources, s.config.Sources)
	return sources
}

// GetSource returns a source by ID
func (s *ConfigService) GetSource(id string) (*model.SourceConfiguration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, source := range s.config.Sources {
		if source.ID == id {
			sourceCopy := source
			return &sourceCopy, nil
		}
	}

	return nil, errors.ErrNotFound
}

// AddSource adds a new source configuration intelligently
// It checks for duplicates and validates the configuration
func (s *ConfigService) AddSource(ctx context.Context, source *model.SourceConfiguration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate source
	if err := s.validateSource(source); err != nil {
		return err
	}

	// Check for duplicate ID
	for _, existing := range s.config.Sources {
		if existing.ID == source.ID {
			return fmt.Errorf("source with ID %s already exists", source.ID)
		}
	}

	// Check for duplicate configuration (intelligent duplicate detection)
	if duplicate := s.findDuplicateSource(source); duplicate != nil {
		return fmt.Errorf("source with similar configuration already exists: %s", duplicate.Name)
	}

	// Set timestamps
	now := time.Now()
	source.CreatedAt = now
	source.UpdatedAt = now

	// Add source
	s.config.Sources = append(s.config.Sources, *source)

	// Save configuration
	if err := s.repo.Save(ctx, s.config); err != nil {
		// Rollback on save failure
		s.config.Sources = s.config.Sources[:len(s.config.Sources)-1]
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// UpdateSource updates an existing source configuration
func (s *ConfigService) UpdateSource(ctx context.Context, source *model.SourceConfiguration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate source
	if err := s.validateSource(source); err != nil {
		return err
	}

	// Find and update source
	found := false
	var oldSource model.SourceConfiguration
	for i, existing := range s.config.Sources {
		if existing.ID == source.ID {
			oldSource = existing
			source.CreatedAt = existing.CreatedAt // Preserve creation time
			source.UpdatedAt = time.Now()
			s.config.Sources[i] = *source
			found = true
			break
		}
	}

	if !found {
		return errors.ErrNotFound
	}

	// Save configuration
	if err := s.repo.Save(ctx, s.config); err != nil {
		// Rollback on save failure
		for i, existing := range s.config.Sources {
			if existing.ID == source.ID {
				s.config.Sources[i] = oldSource
				break
			}
		}
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// RemoveSource removes a source configuration
func (s *ConfigService) RemoveSource(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find source index
	index := -1
	var removedSource model.SourceConfiguration
	for i, source := range s.config.Sources {
		if source.ID == id {
			index = i
			removedSource = source
			break
		}
	}

	if index == -1 {
		return errors.ErrNotFound
	}

	// Remove source
	s.config.Sources = append(s.config.Sources[:index], s.config.Sources[index+1:]...)

	// Save configuration
	if err := s.repo.Save(ctx, s.config); err != nil {
		// Rollback on save failure
		s.config.Sources = append(s.config.Sources[:index], append([]model.SourceConfiguration{removedSource}, s.config.Sources[index:]...)...)
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// validateSource validates a source configuration
func (s *ConfigService) validateSource(source *model.SourceConfiguration) error {
	if source.ID == "" {
		return fmt.Errorf("source ID cannot be empty")
	}

	if source.Name == "" {
		return fmt.Errorf("source name cannot be empty")
	}

	if source.Type == "" {
		return fmt.Errorf("source type cannot be empty")
	}

	// Type-specific validation
	switch source.Type {
	case model.SourceTypeFilesystem:
		return s.validateFilesystemSource(source)
	case model.SourceTypeAPISelfHosted:
		return s.validateAPISource(source)
	default:
		return fmt.Errorf("unsupported source type: %s", source.Type)
	}
}

// validateFilesystemSource validates filesystem-specific configuration
func (s *ConfigService) validateFilesystemSource(source *model.SourceConfiguration) error {
	// Check for root_paths or root_path
	hasRootPaths := false
	hasRootPath := false

	if rootPaths, ok := source.Config["root_paths"].([]interface{}); ok {
		if len(rootPaths) > 0 {
			hasRootPaths = true
		}
	}

	if rootPath, ok := source.Config["root_path"].(string); ok {
		if rootPath != "" {
			hasRootPath = true
		}
	}

	if !hasRootPaths && !hasRootPath {
		return fmt.Errorf("filesystem source must have root_path or root_paths configured")
	}

	return nil
}

// validateAPISource validates API-specific configuration
func (s *ConfigService) validateAPISource(source *model.SourceConfiguration) error {
	if baseURL, ok := source.Config["base_url"].(string); !ok || baseURL == "" {
		return fmt.Errorf("API source must have base_url configured")
	}

	return nil
}

// findDuplicateSource intelligently finds duplicate sources
// For filesystem sources, it checks if root paths overlap
func (s *ConfigService) findDuplicateSource(newSource *model.SourceConfiguration) *model.SourceConfiguration {
	for _, existing := range s.config.Sources {
		if existing.Type != newSource.Type {
			continue
		}

		switch newSource.Type {
		case model.SourceTypeFilesystem:
			if s.isFilesystemDuplicate(&existing, newSource) {
				return &existing
			}
		case model.SourceTypeAPISelfHosted:
			if s.isAPIDuplicate(&existing, newSource) {
				return &existing
			}
		}
	}

	return nil
}

// isFilesystemDuplicate checks if two filesystem sources have the same root path
func (s *ConfigService) isFilesystemDuplicate(existing, newSource *model.SourceConfiguration) bool {
	existingPaths := s.extractRootPaths(existing)
	newPaths := s.extractRootPaths(newSource)

	// Check for exact matches or overlapping paths
	for _, newPath := range newPaths {
		for _, existingPath := range existingPaths {
			if s.pathsOverlap(existingPath, newPath) {
				return true
			}
		}
	}

	return false
}

// isAPIDuplicate checks if two API sources have the same base URL
func (s *ConfigService) isAPIDuplicate(existing, newSource *model.SourceConfiguration) bool {
	existingURL, _ := existing.Config["base_url"].(string)
	newURL, _ := newSource.Config["base_url"].(string)

	return strings.TrimSuffix(existingURL, "/") == strings.TrimSuffix(newURL, "/")
}

// extractRootPaths extracts all root paths from a source configuration
func (s *ConfigService) extractRootPaths(source *model.SourceConfiguration) []string {
	var paths []string

	// Try root_paths array
	if rootPaths, ok := source.Config["root_paths"].([]interface{}); ok {
		for _, p := range rootPaths {
			if path, ok := p.(string); ok {
				paths = append(paths, path)
			}
		}
	}

	// Try root_path string
	if rootPath, ok := source.Config["root_path"].(string); ok && rootPath != "" {
		paths = append(paths, rootPath)
	}

	return paths
}

// pathsOverlap checks if two filesystem paths overlap
func (s *ConfigService) pathsOverlap(path1, path2 string) bool {
	// Normalize paths (remove trailing slashes)
	path1 = strings.TrimSuffix(path1, "/")
	path2 = strings.TrimSuffix(path2, "/")

	// Exact match
	if path1 == path2 {
		return true
	}

	// Check if one is a subdirectory of the other
	if strings.HasPrefix(path1, path2+"/") || strings.HasPrefix(path2, path1+"/") {
		return true
	}

	return false
}