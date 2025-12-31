package model

import "time"

// AppConfig represents the complete application configuration
type AppConfig struct {
	Version string                `json:"version"`
	Sources []SourceConfiguration `json:"sources"`
}

// SourceConfiguration represents a configured music source
type SourceConfiguration struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      SourceType             `json:"type"`
	Enabled   bool                   `json:"enabled"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// FilesystemSourceConfigDTO is the DTO for filesystem source configuration
type FilesystemSourceConfigDTO struct {
	RootPaths        []string `json:"root_paths"`
	IncludeSubfolders bool    `json:"include_subfolders"`
	WatchForChanges  bool     `json:"watch_for_changes"`
	SupportedFormats []string `json:"supported_formats"`
}

// ToFilesystemConfig converts the generic config map to FilesystemSourceConfig
func (sc *SourceConfiguration) ToFilesystemConfig() (*FilesystemSourceConfig, error) {
	// Extract root paths (supports both single and multiple)
	rootPath := ""
	if rootPaths, ok := sc.Config["root_paths"].([]interface{}); ok && len(rootPaths) > 0 {
		if path, ok := rootPaths[0].(string); ok {
			rootPath = path
		}
	} else if path, ok := sc.Config["root_path"].(string); ok {
		rootPath = path
	}

	// Extract watch for changes
	watchForChanges := false
	if watch, ok := sc.Config["watch_for_changes"].(bool); ok {
		watchForChanges = watch
	}

	// Extract supported formats
	var supportedFormats []string
	if formats, ok := sc.Config["supported_formats"].([]interface{}); ok {
		for _, f := range formats {
			if format, ok := f.(string); ok {
				supportedFormats = append(supportedFormats, format)
			}
		}
	}

	return &FilesystemSourceConfig{
		RootPath:         rootPath,
		WatchForChanges:  watchForChanges,
		SupportedFormats: supportedFormats,
	}, nil
}

// NewSourceConfiguration creates a new source configuration
func NewSourceConfiguration(id, name string, sourceType SourceType) *SourceConfiguration {
	now := time.Now()
	return &SourceConfiguration{
		ID:        id,
		Name:      name,
		Type:      sourceType,
		Enabled:   true,
		Config:    make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}