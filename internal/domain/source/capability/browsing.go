package capability

// DirectoryBrowser is a source capability that enables directory browsing
// Sources implementing this interface can list files and directories in their storage
type DirectoryBrowser interface {
	ListDirectory(relativePath string) ([]*FileNode, error)
	GetRootPath() string
}

// FileNode represents a file or directory in a browsable source
type FileNode struct {
	Name        string
	Path        string
	IsDirectory bool
	Size        int64
	Extension   string
}
