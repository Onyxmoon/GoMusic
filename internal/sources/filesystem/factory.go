package filesystem

import (
	"path/filepath"
	"strings"
)

// ExtractorFactory creates extractors based on file type
type ExtractorFactory struct {
	defaultSourceID string
	extractors      map[string]Extractor
}

// NewExtractorFactory creates a new extractor factory
func NewExtractorFactory(defaultSourceID string) *ExtractorFactory {
	factory := &ExtractorFactory{
		defaultSourceID: defaultSourceID,
		extractors:      make(map[string]Extractor),
	}

	// Register the tag extractor for all supported formats
	tagExtractor := NewTagExtractor(defaultSourceID)
	supportedExts := []string{".mp3", ".m4a", ".m4b", ".m4p", ".flac", ".ogg", ".oga"}
	for _, ext := range supportedExts {
		factory.Register(ext, tagExtractor)
	}

	return factory
}

// Register registers an extractor for a file extension
func (f *ExtractorFactory) Register(extension string, extractor Extractor) {
	ext := strings.ToLower(extension)
	f.extractors[ext] = extractor
}

// GetExtractor returns an extractor for the given file path
func (f *ExtractorFactory) GetExtractor(filePath string) Extractor {
	ext := strings.ToLower(filepath.Ext(filePath))
	if extractor, exists := f.extractors[ext]; exists {
		return extractor
	}
	return nil
}

// GetDefaultExtractor returns the default tag extractor
func (f *ExtractorFactory) GetDefaultExtractor() Extractor {
	return NewTagExtractor(f.defaultSourceID)
}

// SupportsFile checks if a file is supported by any registered extractor
func (f *ExtractorFactory) SupportsFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	_, exists := f.extractors[ext]
	return exists
}

// SupportedFormats returns a list of all supported file extensions
func (f *ExtractorFactory) SupportedFormats() []string {
	formats := make([]string, 0, len(f.extractors))
	for ext := range f.extractors {
		formats = append(formats, ext)
	}
	return formats
}