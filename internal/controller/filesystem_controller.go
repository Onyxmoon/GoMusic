package controller

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"GoMusic/internal/application/dto"
	"GoMusic/internal/domain/source/capability"
	"GoMusic/internal/service"
)

// FilesystemController handles filesystem-specific operations like directory browsing
type FilesystemController struct {
	libraryService *service.LibraryService
	ctx            context.Context
}

// NewFilesystemController creates a new FilesystemController
func NewFilesystemController(libraryService *service.LibraryService, ctx context.Context) *FilesystemController {
	return &FilesystemController{
		libraryService: libraryService,
		ctx:            ctx,
	}
}

// BrowseDirectory lists the contents of a directory for a given source
// Returns both files and subdirectories
func (c *FilesystemController) BrowseDirectory(sourceID string, relativePath string) (*dto.DirectoryContentsDTO, error) {
	// Get directory browser for this source
	dirBrowser, err := c.getDirectoryBrowser(sourceID)
	if err != nil {
		return nil, err
	}

	// Get directory contents
	nodes, err := dirBrowser.ListDirectory(relativePath)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs and separate files from directories
	files, directories := convertFileNodesToDTOs(nodes)

	return &dto.DirectoryContentsDTO{
		CurrentPath: relativePath,
		Files:       files,
		Directories: directories,
	}, nil
}

// GetSourceRootPath returns the root path for a filesystem source
func (c *FilesystemController) GetSourceRootPath(sourceID string) (string, error) {
	dirBrowser, err := c.getDirectoryBrowser(sourceID)
	if err != nil {
		return "", err
	}

	return dirBrowser.GetRootPath(), nil
}

// SelectDirectory opens a directory picker dialog
func (c *FilesystemController) SelectDirectory() (string, error) {
	selectedPath, err := runtime.OpenDirectoryDialog(c.ctx, runtime.OpenDialogOptions{
		Title: "Select Music Folder",
	})
	if err != nil {
		return "", err
	}
	return selectedPath, nil
}

// getDirectoryBrowser retrieves a directory browser for a source
func (c *FilesystemController) getDirectoryBrowser(sourceID string) (capability.DirectoryBrowser, error) {
	repos := c.libraryService.GetRepositories()
	repo, ok := repos[sourceID]
	if !ok {
		return nil, fmt.Errorf("source not found: %s", sourceID)
	}

	dirBrowser, ok := repo.(capability.DirectoryBrowser)
	if !ok {
		return nil, fmt.Errorf("source does not support directory browsing")
	}

	return dirBrowser, nil
}

// convertFileNodesToDTOs converts file nodes to DTOs and separates files from directories
func convertFileNodesToDTOs(nodes []*capability.FileNode) (files []*dto.FileNodeDTO, directories []*dto.FileNodeDTO) {
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

	return files, directories
}