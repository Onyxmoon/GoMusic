package dto

// FileNodeDTO represents a file or directory in the filesystem
type FileNodeDTO struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"isDirectory"`
	Size        int64  `json:"size,omitempty"`
	Extension   string `json:"extension,omitempty"`
}

// DirectoryContentsDTO represents the contents of a directory
type DirectoryContentsDTO struct {
	CurrentPath string         `json:"currentPath"`
	Files       []*FileNodeDTO `json:"files"`
	Directories []*FileNodeDTO `json:"directories"`
}