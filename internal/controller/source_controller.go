package controller

import (
	"context"
	"fmt"
	"time"

	"GoMusic/internal/application/dto"
	"GoMusic/internal/domain/model"
	"GoMusic/internal/service"
	"GoMusic/internal/sources/filesystem"
)

// SourceController handles all source management operations
type SourceController struct {
	configService  *service.ConfigService
	libraryService *service.LibraryService
}

// NewSourceController creates a new SourceController
func NewSourceController(configService *service.ConfigService, libraryService *service.LibraryService) *SourceController {
	return &SourceController{
		configService:  configService,
		libraryService: libraryService,
	}
}

// GetSources returns information about all registered sources
func (c *SourceController) GetSources() []dto.SourceDTO {
	sources := c.configService.GetSources()
	result := make([]dto.SourceDTO, len(sources))

	for i, source := range sources {
		result[i] = dto.SourceDTO{
			ID:   source.ID,
			Type: string(source.Type),
			Name: source.Name,
		}
	}

	return result
}

// GetSourceConfig returns the full configuration for a source
func (c *SourceController) GetSourceConfig(sourceID string) (*model.SourceConfiguration, error) {
	source, err := c.configService.GetSource(sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source config: %w", err)
	}
	return source, nil
}

// GetSupportedFormats returns all supported audio file formats
func (c *SourceController) GetSupportedFormats() []string {
	return []string{".mp3", ".flac", ".m4a", ".ogg", ".oga"}
}

// AddFilesystemSource adds a new filesystem music source
func (c *SourceController) AddFilesystemSource(ctx context.Context, name string, rootPaths []string, includeSubfolders bool, formats []string) error {
	// Generate unique source ID based on timestamp
	sourceID := fmt.Sprintf("filesystem-%d", time.Now().Unix())

	// Create source configuration
	sourceConfig := model.NewSourceConfiguration(sourceID, name, model.SourceTypeFilesystem)
	sourceConfig.Config["root_paths"] = convertToInterfaceSlice(rootPaths)
	sourceConfig.Config["root_path"] = rootPaths[0] // For backwards compatibility
	sourceConfig.Config["include_subfolders"] = includeSubfolders
	sourceConfig.Config["watch_for_changes"] = false
	sourceConfig.Config["supported_formats"] = convertToInterfaceSlice(formats)

	// Add to config service (this validates and checks for duplicates)
	if err := c.configService.AddSource(ctx, sourceConfig); err != nil {
		return fmt.Errorf("failed to add source to config: %w", err)
	}

	// Register with library service
	if err := c.registerSource(sourceConfig); err != nil {
		// Rollback config change if registration fails
		_ = c.configService.RemoveSource(ctx, sourceID)
		return fmt.Errorf("failed to register source: %w", err)
	}

	return nil
}

// UpdateFilesystemSource updates an existing filesystem source
func (c *SourceController) UpdateFilesystemSource(ctx context.Context, sourceID string, name string, rootPaths []string, includeSubfolders bool, formats []string) error {
	// Get existing source
	existingSource, err := c.configService.GetSource(sourceID)
	if err != nil {
		return fmt.Errorf("source not found: %w", err)
	}

	// Update configuration
	existingSource.Name = name
	existingSource.Config["root_paths"] = convertToInterfaceSlice(rootPaths)
	existingSource.Config["root_path"] = rootPaths[0]
	existingSource.Config["include_subfolders"] = includeSubfolders
	existingSource.Config["supported_formats"] = convertToInterfaceSlice(formats)

	// Update in config service
	if err := c.configService.UpdateSource(ctx, existingSource); err != nil {
		return fmt.Errorf("failed to update source config: %w", err)
	}

	// Unregister old source from library service
	c.libraryService.UnregisterTrackRepository(sourceID)

	// Re-register with new configuration
	if err := c.registerSource(existingSource); err != nil {
		return fmt.Errorf("failed to re-register source: %w", err)
	}

	return nil
}

// RemoveSource removes a music source
func (c *SourceController) RemoveSource(ctx context.Context, sourceID string) error {
	// Remove from config service (this persists the change)
	if err := c.configService.RemoveSource(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to remove source from config: %w", err)
	}

	// Unregister from library service
	c.libraryService.UnregisterTrackRepository(sourceID)

	return nil
}

// LoadSourcesFromConfig loads all configured sources
func (c *SourceController) LoadSourcesFromConfig() error {
	sources := c.configService.GetSources()

	for _, sourceConfig := range sources {
		if !sourceConfig.Enabled {
			continue
		}

		if err := c.registerSource(&sourceConfig); err != nil {
			fmt.Printf("Failed to register source %s: %v\n", sourceConfig.ID, err)
			continue
		}
	}

	return nil
}

// registerSource registers a source with the library service
func (c *SourceController) registerSource(sourceConfig *model.SourceConfiguration) error {
	switch sourceConfig.Type {
	case model.SourceTypeFilesystem:
		return c.registerFilesystemSource(sourceConfig)
	case model.SourceTypeAPISelfHosted:
		// TODO: Implement API source registration
		return fmt.Errorf("API sources not yet implemented")
	default:
		return fmt.Errorf("unsupported source type: %s", sourceConfig.Type)
	}
}

// registerFilesystemSource registers a filesystem source
func (c *SourceController) registerFilesystemSource(sourceConfig *model.SourceConfiguration) error {
	// Convert config to FilesystemSourceConfig
	config, err := sourceConfig.ToFilesystemConfig()
	if err != nil {
		return fmt.Errorf("failed to convert config: %w", err)
	}

	// Create extractor (filesystem-specific)
	extractor := filesystem.NewTagExtractor(sourceConfig.ID)

	// Create repository
	repo := filesystem.NewFilesystemTrackRepository(sourceConfig.ID, config, extractor)

	// Register with library service
	c.libraryService.RegisterTrackRepository(sourceConfig.ID, repo)

	return nil
}

// convertToInterfaceSlice converts a string slice to interface slice for JSON marshaling
func convertToInterfaceSlice(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, s := range strings {
		interfaces[i] = s
	}
	return interfaces
}