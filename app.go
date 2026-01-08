package main

import (
	"GoMusic/internal/domain/model"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"GoMusic/internal/application/dto"
	"GoMusic/internal/application/mapper"
	"GoMusic/internal/controller"
	"GoMusic/internal/domain/repository"
	configRepo "GoMusic/internal/repository/config"
	"GoMusic/internal/service"
)

// App struct
type App struct {
	ctx context.Context

	// Services
	libraryService *service.LibraryService
	configService  *service.ConfigService

	// Controllers
	sourceController     *controller.SourceController
	scanController       *controller.ScanController
	filesystemController *controller.FilesystemController

	// Mappers
	trackMapper *mapper.TrackMapper
}

// NewApp creates a new App application struct
func NewApp() *App {
	libraryService := service.NewLibraryService()

	// Initialize config service with JSON repository
	configPath := getConfigPath()
	configRepository := configRepo.NewJSONConfigRepository(configPath)
	configService := service.NewConfigService(configRepository)

	return &App{
		libraryService: libraryService,
		configService:  configService,
		trackMapper:    mapper.NewTrackMapper(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize controllers
	a.sourceController = controller.NewSourceController(a.configService, a.libraryService)
	a.scanController = controller.NewScanController(a.libraryService, ctx)
	a.filesystemController = controller.NewFilesystemController(a.libraryService, ctx)

	// Initialize configuration
	if err := a.configService.Initialize(ctx); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		return
	}

	// Load sources from configuration
	if err := a.sourceController.LoadSourcesFromConfig(); err != nil {
		fmt.Printf("Failed to load sources from config: %v\n", err)
	}
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "./gomusic-config.json"
	}

	// Use ~/.gomusic/config.json
	configDir := filepath.Join(homeDir, ".gomusic")
	return filepath.Join(configDir, "config.json")
}

// === WAILS-EXPOSED METHODS (callable from Svelte frontend) ===

// GetAllTracks returns all tracks from all sources
func (a *App) GetAllTracks() ([]*dto.TrackDTO, error) {
	tracks, err := a.libraryService.GetAllTracks(a.ctx, nil)
	if err != nil {
		return nil, err
	}

	return a.trackMapper.ToDTOList(tracks), nil
}

// GetTrack retrieves a single track by ID
func (a *App) GetTrack(id string) (*dto.TrackDTO, error) {
	track, err := a.libraryService.GetTrackByID(a.ctx, id)
	if err != nil {
		return nil, err
	}

	return a.trackMapper.ToDTO(track), nil
}

// SearchTracks searches for tracks across all sources
func (a *App) SearchTracks(query string) ([]*dto.TrackDTO, error) {
	opts := repository.DefaultSearchOptions()
	tracks, err := a.libraryService.SearchTracks(a.ctx, query, opts)
	if err != nil {
		return nil, err
	}

	return a.trackMapper.ToDTOList(tracks), nil
}

// GetTracksByAlbum retrieves all tracks for a specific album
func (a *App) GetTracksByAlbum(albumID string) ([]*dto.TrackDTO, error) {
	tracks, err := a.libraryService.GetTracksByAlbum(a.ctx, albumID)
	if err != nil {
		return nil, err
	}

	return a.trackMapper.ToDTOList(tracks), nil
}

// GetTracksByArtist retrieves all tracks for a specific artist
func (a *App) GetTracksByArtist(artistID string) ([]*dto.TrackDTO, error) {
	tracks, err := a.libraryService.GetTracksByArtist(a.ctx, artistID)
	if err != nil {
		return nil, err
	}

	return a.trackMapper.ToDTOList(tracks), nil
}

// === Scan Operations (delegated to ScanController) ===

// ScanLibrary triggers a library scan for a specific source
func (a *App) ScanLibrary(sourceID string) error {
	return a.scanController.ScanLibrary(sourceID)
}

// ScanAllLibraries triggers a scan on all registered sources
func (a *App) ScanAllLibraries() error {
	return a.scanController.ScanAllLibraries()
}

// GetScanProgress retrieves current scan progress for a source
func (a *App) GetScanProgress(sourceID string) (*dto.ScanProgressDTO, error) {
	return a.scanController.GetScanProgress(sourceID)
}

// GetAllScanProgress retrieves scan progress for all sources
func (a *App) GetAllScanProgress() map[string]*dto.ScanProgressDTO {
	return a.scanController.GetAllScanProgress()
}

// === Source Management (delegated to SourceController) ===

// GetSources returns information about all registered sources
func (a *App) GetSources() []dto.SourceDTO {
	return a.sourceController.GetSources()
}

// GetSourceConfig returns the full configuration for a source
func (a *App) GetSourceConfig(sourceID string) (*model.SourceConfiguration, error) {
	return a.sourceController.GetSourceConfig(sourceID)
}

// GetSupportedFormats returns all supported audio file formats
func (a *App) GetSupportedFormats() []string {
	return a.sourceController.GetSupportedFormats()
}

// AddFilesystemSource adds a new filesystem music source
func (a *App) AddFilesystemSource(name string, rootPaths []string, includeSubfolders bool, formats []string) error {
	return a.sourceController.AddFilesystemSource(a.ctx, name, rootPaths, includeSubfolders, formats)
}

// UpdateFilesystemSource updates an existing filesystem source
func (a *App) UpdateFilesystemSource(sourceID string, name string, rootPaths []string, includeSubfolders bool, formats []string) error {
	return a.sourceController.UpdateFilesystemSource(a.ctx, sourceID, name, rootPaths, includeSubfolders, formats)
}

// RemoveSource removes a music source
func (a *App) RemoveSource(sourceID string) error {
	return a.sourceController.RemoveSource(a.ctx, sourceID)
}

// === Filesystem Operations (delegated to FilesystemController) ===

// BrowseDirectory lists the contents of a directory for a given source
func (a *App) BrowseDirectory(sourceID string, relativePath string) (*dto.DirectoryContentsDTO, error) {
	return a.filesystemController.BrowseDirectory(sourceID, relativePath)
}

// GetSourceRootPath returns the root path for a filesystem source
func (a *App) GetSourceRootPath(sourceID string) (string, error) {
	return a.filesystemController.GetSourceRootPath(sourceID)
}

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	return a.filesystemController.SelectDirectory()
}

// === HTTP MIDDLEWARE ===

// AudioFileMiddleware intercepts audio streaming and artwork requests
func (a *App) AudioFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle requests to /audio/* path
		if len(r.URL.Path) >= 7 && r.URL.Path[:7] == "/audio/" {
			a.serveAudioFile(w, r)
			return
		}

		// Handle requests to /artwork/* path
		if len(r.URL.Path) >= 9 && r.URL.Path[:9] == "/artwork/" {
			a.serveArtworkFile(w, r)
			return
		}

		// Not an audio or artwork request, pass to next handler
		next.ServeHTTP(w, r)
	})
}

// === AUDIO PLAYBACK METHODS ===

// GetTrackFilePath returns the file path for a track
// This works for local filesystem sources
// Later: extend to return stream URLs for remote sources
func (a *App) GetTrackFilePath(trackID string) (string, error) {
	track, err := a.libraryService.GetTrackByID(a.ctx, trackID)
	if err != nil {
		return "", fmt.Errorf("track not found: %w", err)
	}

	// For now: return FilePath for filesystem sources
	// Later: check track.SourceType and return appropriate URL
	//   - filesystem → file:// path
	//   - api → http:// stream URL
	//   - remote → https:// stream URL
	return track.FilePath, nil
}

// serveAudioFile handles HTTP requests for streaming audio files by track ID
func (a *App) serveAudioFile(w http.ResponseWriter, r *http.Request) {
	// Extract track ID from URL query parameter
	trackID := r.URL.Query().Get("id")
	if trackID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	// Look up track in library service cache
	track, err := a.libraryService.GetTrackByID(a.ctx, trackID)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	// Get file path from track metadata
	filePath := track.FilePath
	if filePath == "" {
		http.Error(w, "Track has no file path", http.StatusInternalServerError)
		return
	}

	// Verify file exists and is readable
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Cannot access file", http.StatusInternalServerError)
		}
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, "Path is a directory, not a file", http.StatusBadRequest)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusNotFound)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	// Set response headers
	w.Header().Set("Content-Type", getAudioContentType(filePath))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Stream the file
	_, _ = io.Copy(w, file)
}

// serveArtworkFile handles HTTP requests for serving album artwork by track ID
func (a *App) serveArtworkFile(w http.ResponseWriter, r *http.Request) {
	// Extract track ID from URL query parameter
	trackID := r.URL.Query().Get("id")
	if trackID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	// Look up track in library service cache
	track, err := a.libraryService.GetTrackByID(a.ctx, trackID)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	// Get artwork path from track metadata
	artworkFilename := track.ArtworkPath
	if artworkFilename == "" {
		http.Error(w, "Track has no artwork", http.StatusNotFound)
		return
	}

	// Construct full path to artwork file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
		return
	}
	artworkDir := filepath.Join(homeDir, ".gomusic", "artwork")
	fullPath := filepath.Join(artworkDir, artworkFilename)

	// Open the artwork file
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "Artwork not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get file info for Content-Length
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Cannot stat file", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", getImageContentType(fullPath))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Stream the image
	if _, err := io.Copy(w, file); err != nil {
		fmt.Printf("Error streaming artwork: %v\n", err)
	}
}

// getAudioContentType determines the MIME type for audio files
func getAudioContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".flac":
		return "audio/flac"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}

// getImageContentType determines the MIME type for image files
func getImageContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
