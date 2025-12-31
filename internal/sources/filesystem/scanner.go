package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"GoMusic/internal/domain/source/capability"
)

// DirectoryScanner scans directories for audio files
type DirectoryScanner struct {
	rootPath         string
	supportedFormats map[string]bool
	ignoredPrefixes  []string 
	ignoredFiles     []string
}

// NewDirectoryScanner creates a new directory scanner
func NewDirectoryScanner(rootPath string, supportedFormats []string) *DirectoryScanner {
	formats := make(map[string]bool)
	for _, format := range supportedFormats {
		formats[strings.ToLower(format)] = true
	}

	return &DirectoryScanner{
		rootPath:         rootPath,
		supportedFormats: formats,
		ignoredPrefixes:  []string{"._"}, // macOS AppleDouble files
		ignoredFiles:     []string{".DS_Store", "Thumbs.db"}, // macOS and Windows metadata
	}
}

// ScanDirectory recursively scans a directory for audio files
// The callback is called for each file found, useful for progress tracking
func (s *DirectoryScanner) ScanDirectory(ctx context.Context, callback func(filePath string)) ([]string, error) {
	var audioFiles []string

	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		// Check for context cancellation (user cancelled scan or application is shutting down)
		select {
		case <-ctx.Done():

			return ctx.Err()
		default:
		}

		if err != nil {
			// Log error but continue walking
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip ignored files
		fileName := filepath.Base(path)
		if s.shouldIgnoreFile(fileName) {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if s.supportedFormats[ext] {
			audioFiles = append(audioFiles, path)
			if callback != nil {
				callback(path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return audioFiles, nil
}

// IsSupported checks if a file extension is supported
func (s *DirectoryScanner) IsSupported(extension string) bool {
	return s.supportedFormats[strings.ToLower(extension)]
}

// shouldIgnoreFile checks if a file should be ignored based on ignore lists
func (s *DirectoryScanner) shouldIgnoreFile(fileName string) bool {
	// Check ignored file names
	for _, ignoredFile := range s.ignoredFiles {
		if fileName == ignoredFile {
			return true
		}
	}

	// Check ignored prefixes
	for _, prefix := range s.ignoredPrefixes {
		if strings.HasPrefix(fileName, prefix) {
			return true
		}
	}

	return false
}

// ListDirectory lists the immediate contents of a directory (non-recursive)
// Returns both files and subdirectories
func (s *DirectoryScanner) ListDirectory(dirPath string) ([]*capability.FileNode, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var nodes []*capability.FileNode

	for _, entry := range entries {
		// Skip ignored files
		if s.shouldIgnoreFile(entry.Name()) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			// Skip entries we can't read
			continue
		}

		node := &capability.FileNode{
			Name:        entry.Name(),
			Path:        filepath.Join(dirPath, entry.Name()),
			IsDirectory: entry.IsDir(),
		}

		if !entry.IsDir() {
			node.Size = info.Size()
			node.Extension = strings.ToLower(filepath.Ext(entry.Name()))
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
