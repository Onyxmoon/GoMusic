package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"GoMusic/internal/domain/model"
)

// JSONConfigRepository implements ConfigRepository using JSON files
type JSONConfigRepository struct {
	configPath string
}

// NewJSONConfigRepository creates a new JSON config repository
func NewJSONConfigRepository(configPath string) *JSONConfigRepository {
	return &JSONConfigRepository{
		configPath: configPath,
	}
}

// Load loads the application configuration from JSON file
func (r *JSONConfigRepository) Load(ctx context.Context) (*model.AppConfig, error) {
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config model.AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves the application configuration to JSON file
func (r *JSONConfigRepository) Save(ctx context.Context, config *model.AppConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if the configuration file exists
func (r *JSONConfigRepository) Exists(ctx context.Context) (bool, error) {
	_, err := os.Stat(r.configPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Initialize creates a new default configuration
func (r *JSONConfigRepository) Initialize(ctx context.Context) (*model.AppConfig, error) {
	config := &model.AppConfig{
		Version: "1.0.0",
		Sources: []model.SourceConfiguration{},
	}

	if err := r.Save(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}