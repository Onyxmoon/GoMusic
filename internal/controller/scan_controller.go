package controller

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"GoMusic/internal/application/dto"
	"GoMusic/internal/service"
)

// ScanController handles all library scanning operations
type ScanController struct {
	libraryService *service.LibraryService
	ctx            context.Context
}

// NewScanController creates a new ScanController
func NewScanController(libraryService *service.LibraryService, ctx context.Context) *ScanController {
	return &ScanController{
		libraryService: libraryService,
		ctx:            ctx,
	}
}

// ScanLibrary triggers a library scan for a specific source
// Runs asynchronously and emits events for progress updates
func (c *ScanController) ScanLibrary(sourceID string) error {
	return c.runScanWithEvents(sourceID, func(ctx context.Context) error {
		return c.libraryService.ScanSource(ctx, sourceID)
	})
}

// ScanAllLibraries triggers a scan on all registered sources
func (c *ScanController) ScanAllLibraries() error {
	return c.runScanWithEvents("all", func(ctx context.Context) error {
		return c.libraryService.ScanAllSources(ctx)
	})
}

// GetScanProgress retrieves current scan progress for a source
func (c *ScanController) GetScanProgress(sourceID string) (*dto.ScanProgressDTO, error) {
	progress, err := c.libraryService.GetScanProgress(sourceID)
	if err != nil {
		return nil, err
	}

	return dto.ToScanProgressDTO(progress), nil
}

// GetAllScanProgress retrieves scan progress for all sources
func (c *ScanController) GetAllScanProgress() map[string]*dto.ScanProgressDTO {
	allProgress := c.libraryService.GetAllScanProgress()
	result := make(map[string]*dto.ScanProgressDTO)

	for sourceID, progress := range allProgress {
		result[sourceID] = dto.ToScanProgressDTO(progress)
	}

	return result
}

// runScanWithEvents is a helper that runs a scan function asynchronously with event emission
func (c *ScanController) runScanWithEvents(sourceID string, scanFunc func(context.Context) error) error {
	go func() {
		ctx, cancel := context.WithCancel(c.ctx)
		defer cancel()

		runtime.EventsEmit(c.ctx, "scan:started", sourceID)

		err := scanFunc(ctx)
		if err != nil {
			runtime.EventsEmit(c.ctx, "scan:error", map[string]interface{}{
				"sourceId": sourceID,
				"error":    err.Error(),
			})
			return
		}

		runtime.EventsEmit(c.ctx, "scan:complete", sourceID)
	}()

	return nil
}