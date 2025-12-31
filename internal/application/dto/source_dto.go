package dto

import "GoMusic/internal/domain/repository"

// SourceDTO is the data transfer object for music sources
type SourceDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ScanProgressDTO is the data transfer object for scan progress
type ScanProgressDTO struct {
	IsScanning     bool     `json:"isScanning"`
	TotalFiles     int      `json:"totalFiles"`
	ProcessedFiles int      `json:"processedFiles"`
	CurrentFile    string   `json:"currentFile"`
	Errors         []string `json:"errors,omitempty"`
}

// ToScanProgressDTO converts repository.ScanProgress to DTO
func ToScanProgressDTO(progress *repository.ScanProgress) *ScanProgressDTO {
	if progress == nil {
		return nil
	}

	return &ScanProgressDTO{
		IsScanning:     progress.IsScanning,
		TotalFiles:     progress.TotalFiles,
		ProcessedFiles: progress.ProcessedFiles,
		CurrentFile:    progress.CurrentFile,
		Errors:         progress.Errors,
	}
}