package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "GoMusic",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: AudioFileMiddleware,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// AudioFileMiddleware intercepts audio streaming and artwork requests
func AudioFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle requests to /audio/* path
		if len(r.URL.Path) >= 7 && r.URL.Path[:7] == "/audio/" {
			serveAudioFile(w, r)
			return
		}

		// Handle requests to /artwork/* path
		if len(r.URL.Path) >= 9 && r.URL.Path[:9] == "/artwork/" {
			serveArtworkFile(w, r)
			return
		}

		// Not an audio or artwork request, pass to next handler
		next.ServeHTTP(w, r)
	})
}

// serveAudioFile handles the actual audio file streaming
func serveAudioFile(w http.ResponseWriter, r *http.Request) {
	// Extract file path from URL query parameter
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	// Security: Validate and clean the file path to prevent directory traversal attacks
	cleanPath := filepath.Clean(filePath)

	// Ensure the path is absolute (prevents relative path attacks)
	if !filepath.IsAbs(cleanPath) {
		http.Error(w, "Invalid file path: must be absolute path", http.StatusBadRequest)
		return
	}

	// Additional security check: reject paths containing ".." segments
	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Invalid file path: path traversal not allowed", http.StatusForbidden)
		return
	}

	// Verify the file exists and is a regular file (not a directory or symlink)
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Cannot access file", http.StatusInternalServerError)
		}
		return
	}

	// Ensure it's a regular file, not a directory
	if fileInfo.IsDir() {
		http.Error(w, "Path is a directory, not a file", http.StatusBadRequest)
		return
	}

	fmt.Printf("🎵 Streaming audio file: %s\n", cleanPath)

	// Open the file
	file, err := os.Open(cleanPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	// Set appropriate headers (reuse fileInfo from earlier Stat call)
	w.Header().Set("Content-Type", getContentType(filePath))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Printf("Streaming %s (%d bytes)\n", getContentType(filePath), fileInfo.Size())

	// Stream the file
	if _, err := io.Copy(w, file); err != nil {
		fmt.Printf("Error streaming file: %v\n", err)
	}
}

// serveArtworkFile handles serving album artwork images
func serveArtworkFile(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL query parameter
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}

	// Construct full path to artwork file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
		return
	}
	artworkDir := filepath.Join(homeDir, ".gomusic", "artwork")
	filePath := filepath.Join(artworkDir, filename)

	// Security: Ensure the resolved path is still within the artwork directory
	// This prevents directory traversal attacks (e.g., file=../../etc/passwd)
	cleanPath := filepath.Clean(filePath)
	cleanArtworkDir := filepath.Clean(artworkDir)

	// Use filepath.Rel to check if the path is within artworkDir
	relPath, err := filepath.Rel(cleanArtworkDir, cleanPath)
	if err != nil || strings.HasPrefix(relPath, "..") || filepath.IsAbs(relPath) {
		http.Error(w, "Invalid file path: access denied", http.StatusForbidden)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
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

	// Set appropriate headers
	w.Header().Set("Content-Type", getImageContentType(filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Stream the image
	if _, err := io.Copy(w, file); err != nil {
		fmt.Printf("Error streaming artwork: %v\n", err)
	}
}

func getContentType(filePath string) string {
	// Determine content type based on file extension
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

func getImageContentType(filePath string) string {
	// Determine content type based on file extension
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
