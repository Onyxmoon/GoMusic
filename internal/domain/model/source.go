package model

import (
	"time"
)

// SourceType defines the type of music source
type SourceType string

const (
	SourceTypeFilesystem    SourceType = "filesystem"
	SourceTypeAPISelfHosted SourceType = "api-selfhosted"
)

// Source represents a music source (filesystem or API)
type Source struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Type    SourceType `json:"type"`
	Enabled bool       `json:"enabled"`

	// Type-specific configuration
	Config SourceConfig `json:"config"`

	// Timestamps
	AddedAt     time.Time `json:"addedAt"`
	LastScanned time.Time `json:"lastScanned,omitempty"`
}

// SourceConfig is an interface for source-specific configuration
type SourceConfig interface {
	Validate() error
	GetType() SourceType
}

// FilesystemSourceConfig holds configuration for filesystem sources
type FilesystemSourceConfig struct {
	RootPath         string   `json:"rootPath"`
	WatchForChanges  bool     `json:"watchForChanges"`
	SupportedFormats []string `json:"supportedFormats"` // [".mp3", ".flac", ".m4a", ".ogg"]
}

// Validate validates the filesystem source configuration
func (c *FilesystemSourceConfig) Validate() error {
	if c.RootPath == "" {
		return ErrInvalidConfig("root path is required")
	}
	if len(c.SupportedFormats) == 0 {
		c.SupportedFormats = []string{".mp3", ".flac", ".m4a", ".ogg"}
	}
	return nil
}

// GetType returns the source type
func (c *FilesystemSourceConfig) GetType() SourceType {
	return SourceTypeFilesystem
}

// APISourceConfig holds configuration for API-based sources
type APISourceConfig struct {
	BaseURL   string        `json:"baseUrl"`
	APIKey    string        `json:"apiKey,omitempty"`
	Timeout   time.Duration `json:"timeout"`
	RateLimit int           `json:"rateLimit"` // Requests per second
}

// Validate validates the API source configuration
func (c *APISourceConfig) Validate() error {
	if c.BaseURL == "" {
		return ErrInvalidConfig("base URL is required")
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	if c.RateLimit == 0 {
		c.RateLimit = 10
	}
	return nil
}

// GetType returns the source type
func (c *APISourceConfig) GetType() SourceType {
	return SourceTypeAPISelfHosted
}

// ErrInvalidConfig creates an invalid configuration error
func ErrInvalidConfig(message string) error {
	return &ConfigError{Message: message}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "invalid configuration: " + e.Message
}
