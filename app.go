package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"GoMusic/internal/application/dto"
	"GoMusic/internal/application/mapper"
	"GoMusic/internal/domain/model"
	"GoMusic/internal/domain/repository"
	"GoMusic/internal/domain/source/capability"
	configRepo "GoMusic/internal/repository/config"
	"GoMusic/internal/service"
	"GoMusic/internal/sources/filesystem"
)

// App struct
type App struct {
	ctx context.Context

	// Services
	libraryService *service.LibraryService
	configService  *service.ConfigService

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

	// Initialize configuration
	if err := a.configService.Initialize(ctx); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		return
	}

	// Load sources from configuration
	if err := a.loadSourcesFromConfig(); err != nil {
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

// loadSourcesFromConfig loads all configured sources
func (a *App) loadSourcesFromConfig() error {
	sources := a.configService.GetSources()

	for _, sourceConfig := range sources {
		if !sourceConfig.Enabled {
			continue
		}

		if err := a.registerSource(&sourceConfig); err != nil {
			fmt.Printf("Failed to register source %s: %v\n", sourceConfig.ID, err)
			continue
		}
	}

	return nil
}

// registerSource registers a source with the library service
func (a *App) registerSource(sourceConfig *model.SourceConfiguration) error {
	switch sourceConfig.Type {
	case model.SourceTypeFilesystem:
		return a.registerFilesystemSource(sourceConfig)
	case model.SourceTypeAPISelfHosted:
		// TODO: Implement API source registration
		return fmt.Errorf("API sources not yet implemented")
	default:
		return fmt.Errorf("unsupported source type: %s", sourceConfig.Type)
	}
}

// registerFilesystemSource registers a filesystem source
func (a *App) registerFilesystemSource(sourceConfig *model.SourceConfiguration) error {
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
	a.libraryService.RegisterTrackRepository(sourceConfig.ID, repo)

	return nil
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

// ScanLibrary triggers a library scan for a specific source
// Runs asynchronously and emits events for progress updates
func (a *App) ScanLibrary(sourceID string) error {
	// Run scan in background to avoid blocking UI
	go func() {
		ctx, cancel := context.WithCancel(a.ctx)
		defer cancel()

		runtime.EventsEmit(a.ctx, "scan:started", sourceID)

		err := a.libraryService.ScanSource(ctx, sourceID)
		if err != nil {
			runtime.EventsEmit(a.ctx, "scan:error", map[string]interface{}{
				"sourceId": sourceID,
				"error":    err.Error(),
			})
			return
		}

		runtime.EventsEmit(a.ctx, "scan:complete", sourceID)
	}()

	return nil
}

// ScanAllLibraries triggers a scan on all registered sources
func (a *App) ScanAllLibraries() error {
	go func() {
		ctx, cancel := context.WithCancel(a.ctx)
		defer cancel()

		runtime.EventsEmit(a.ctx, "scan:started", "all")

		err := a.libraryService.ScanAllSources(ctx)
		if err != nil {
			runtime.EventsEmit(a.ctx, "scan:error", map[string]interface{}{
				"sourceId": "all",
				"error":    err.Error(),
			})
			return
		}

		runtime.EventsEmit(a.ctx, "scan:complete", "all")
	}()

	return nil
}

// GetScanProgress retrieves current scan progress for a source
func (a *App) GetScanProgress(sourceID string) (*dto.ScanProgressDTO, error) {
	progress, err := a.libraryService.GetScanProgress(sourceID)
	if err != nil {
		return nil, err
	}

	return dto.ToScanProgressDTO(progress), nil
}

// GetAllScanProgress retrieves scan progress for all sources
func (a *App) GetAllScanProgress() map[string]*dto.ScanProgressDTO {
	allProgress := a.libraryService.GetAllScanProgress()
	result := make(map[string]*dto.ScanProgressDTO)

	for sourceID, progress := range allProgress {
		result[sourceID] = dto.ToScanProgressDTO(progress)
	}

	return result
}

// GetSources returns information about all registered sources
func (a *App) GetSources() []dto.SourceDTO {
	sources := a.configService.GetSources()
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
func (a *App) GetSourceConfig(sourceID string) (*model.SourceConfiguration, error) {
	source, err := a.configService.GetSource(sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source config: %w", err)
	}
	return source, nil
}

// GetSupportedFormats returns all supported audio file formats
func (a *App) GetSupportedFormats() []string {
	return []string{".mp3", ".flac", ".m4a", ".ogg", ".oga"}
}

// AddFilesystemSource adds a new filesystem music source
func (a *App) AddFilesystemSource(name string, rootPaths []string, includeSubfolders bool, formats []string) error {
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
	if err := a.configService.AddSource(a.ctx, sourceConfig); err != nil {
		return fmt.Errorf("failed to add source to config: %w", err)
	}

	// Register with library service
	if err := a.registerSource(sourceConfig); err != nil {
		// Rollback config change if registration fails
		_ = a.configService.RemoveSource(a.ctx, sourceID)
		return fmt.Errorf("failed to register source: %w", err)
	}

	return nil
}

// UpdateFilesystemSource updates an existing filesystem source
func (a *App) UpdateFilesystemSource(sourceID string, name string, rootPaths []string, includeSubfolders bool, formats []string) error {
	// Get existing source
	existingSource, err := a.configService.GetSource(sourceID)
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
	if err := a.configService.UpdateSource(a.ctx, existingSource); err != nil {
		return fmt.Errorf("failed to update source config: %w", err)
	}

	// Unregister old source from library service
	a.libraryService.UnregisterTrackRepository(sourceID)

	// Re-register with new configuration
	if err := a.registerSource(existingSource); err != nil {
		return fmt.Errorf("failed to re-register source: %w", err)
	}

	return nil
}

// RemoveSource removes a music source
func (a *App) RemoveSource(sourceID string) error {
	// Remove from config service (this persists the change)
	if err := a.configService.RemoveSource(a.ctx, sourceID); err != nil {
		return fmt.Errorf("failed to remove source from config: %w", err)
	}

	// Unregister from library service
	a.libraryService.UnregisterTrackRepository(sourceID)

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

// BrowseDirectory lists the contents of a directory for a given source
// Returns both files and subdirectories
func (a *App) BrowseDirectory(sourceID string, relativePath string) (*dto.DirectoryContentsDTO, error) {
	// Get the repository for this source
	repos := a.libraryService.GetRepositories()
	repo, ok := repos[sourceID]
	if !ok {
		return nil, fmt.Errorf("source not found: %s", sourceID)
	}

	// Check if source supports directory browsing
	dirBrowser, ok := repo.(capability.DirectoryBrowser)
	if !ok {
		return nil, fmt.Errorf("source does not support directory browsing")
	}

	// Get directory contents
	nodes, err := dirBrowser.ListDirectory(relativePath)
	if err != nil {
		return nil, err
	}

	// Separate files and directories
	var files []*dto.FileNodeDTO
	var directories []*dto.FileNodeDTO

	for _, node := range nodes {
		fileNode := &dto.FileNodeDTO{
			Name:        node.Name,
			Path:        node.Path,
			IsDirectory: node.IsDirectory,
			Size:        node.Size,
			Extension:   node.Extension,
		}

		if node.IsDirectory {
			directories = append(directories, fileNode)
		} else {
			files = append(files, fileNode)
		}
	}

	return &dto.DirectoryContentsDTO{
		CurrentPath: relativePath,
		Files:       files,
		Directories: directories,
	}, nil
}

// GetSourceRootPath returns the root path for a filesystem source
func (a *App) GetSourceRootPath(sourceID string) (string, error) {
	repos := a.libraryService.GetRepositories()
	repo, ok := repos[sourceID]
	if !ok {
		return "", fmt.Errorf("source not found: %s", sourceID)
	}

	dirBrowser, ok := repo.(capability.DirectoryBrowser)
	if !ok {
		return "", fmt.Errorf("source is not a filesystem source")
	}

	return dirBrowser.GetRootPath(), nil
}

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	selectedPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Music Folder",
	})
	if err != nil {
		return "", err
	}
	return selectedPath, nil
}

// === HTTP MIDDLEWARE ===

// AudioFileMiddleware intercepts audio streaming and artwork requests
func (a *App) AudioFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle requests to /audio/* path
		if len(r.URL.Path) >= 7 && r.URL.Path[:7] == "/audio/" {
			a.ServeAudioFile(w, r)
			return
		}

		// Handle requests to /artwork/* path
		if len(r.URL.Path) >= 9 && r.URL.Path[:9] == "/artwork/" {
			a.ServeArtworkFile(w, r)
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

// ServeAudioFile handles HTTP requests for streaming audio files by track ID
func (a *App) ServeAudioFile(w http.ResponseWriter, r *http.Request) {
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
	if _, err := io.Copy(w, file); err != nil {
		fmt.Printf("Error streaming audio file: %v\n", err)
	}
}

// ServeArtworkFile handles HTTP requests for serving album artwork by track ID
func (a *App) ServeArtworkFile(w http.ResponseWriter, r *http.Request) {
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

