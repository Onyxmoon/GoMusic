package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/mewkiz/flac"
)

// AudioAnalyzer extracts audio properties from local audio files
// This is specific to filesystem sources - API sources will get this info from the API
type AudioAnalyzer struct{}

// NewAudioAnalyzer creates a new audio analyzer
func NewAudioAnalyzer() *AudioAnalyzer {
	return &AudioAnalyzer{}
}

// AudioProperties contains extracted audio properties
type AudioProperties struct {
	Duration   time.Duration
	SampleRate int
	BitRate    int
}

// Analyze extracts audio properties from a file
func (a *AudioAnalyzer) Analyze(filePath string) *AudioProperties {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".mp3":
		return a.analyzeMP3(filePath)
	case ".flac":
		return a.analyzeFLAC(filePath)
	default:
		// For other formats (M4A, OGG), return empty properties
		// Can be extended later with additional libraries
		return &AudioProperties{}
	}
}

// analyzeMP3 extracts properties from MP3 file
func (a *AudioAnalyzer) analyzeMP3(filePath string) *AudioProperties {
	file, err := os.Open(filePath)
	if err != nil {
		return &AudioProperties{}
	}
	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return &AudioProperties{}
	}

	props := &AudioProperties{}

	// Sample rate
	sampleRate := decoder.SampleRate()
	props.SampleRate = sampleRate

	// Calculate duration
	if sampleRate > 0 {
		length := decoder.Length()
		// Length is in bytes (4 bytes per sample for stereo 16-bit)
		samples := length / 4
		durationSeconds := float64(samples) / float64(sampleRate)
		props.Duration = time.Duration(durationSeconds * float64(time.Second))
	}

	// Estimate bitrate from file size and duration
	if props.Duration > 0 {
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			fileSizeBytes := fileInfo.Size()
			durationSeconds := props.Duration.Seconds()
			props.BitRate = int((float64(fileSizeBytes) * 8) / (durationSeconds * 1000))
		}
	}

	return props
}

// analyzeFLAC extracts properties from FLAC file
func (a *AudioAnalyzer) analyzeFLAC(filePath string) *AudioProperties {
	stream, err := flac.ParseFile(filePath)
	if err != nil {
		return &AudioProperties{}
	}

	info := stream.Info
	props := &AudioProperties{}

	// Sample rate
	props.SampleRate = int(info.SampleRate)

	// Duration
	if info.SampleRate > 0 {
		durationSeconds := float64(info.NSamples) / float64(info.SampleRate)
		props.Duration = time.Duration(durationSeconds * float64(time.Second))
	}

	// Estimate bitrate
	if props.Duration > 0 {
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			fileSizeBytes := fileInfo.Size()
			durationSeconds := props.Duration.Seconds()
			props.BitRate = int((float64(fileSizeBytes) * 8) / (durationSeconds * 1000))
		}
	}

	return props
}