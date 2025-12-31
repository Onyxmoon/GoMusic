package errors

import "errors"

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrSourceNotFound is returned when a source is not found
	ErrSourceNotFound = errors.New("source not found")

	// ErrScanInProgress is returned when a scan is already in progress
	ErrScanInProgress = errors.New("scan already in progress")

	// ErrUnsupportedFormat is returned when a file format is not supported
	ErrUnsupportedFormat = errors.New("unsupported file format")

	// ErrMetadataExtraction is returned when metadata extraction fails
	ErrMetadataExtraction = errors.New("failed to extract metadata")

	// ErrInvalidConfiguration is returned when configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")
)

// NotFoundError creates a not found error with a custom message
func NotFoundError(resource string) error {
	return &ResourceError{
		Type:    "not_found",
		Message: resource + " not found",
	}
}

// ValidationError creates a validation error with a custom message
func ValidationError(field, message string) error {
	return &ResourceError{
		Type:    "validation",
		Message: field + ": " + message,
	}
}

// ResourceError represents a resource-specific error
type ResourceError struct {
	Type    string
	Message string
}

func (e *ResourceError) Error() string {
	return e.Message
}

// Is checks if the error is of a specific type
func (e *ResourceError) Is(target error) bool {
	if t, ok := target.(*ResourceError); ok {
		return e.Type == t.Type
	}
	return false
}