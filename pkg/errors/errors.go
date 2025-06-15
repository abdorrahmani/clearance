package errors

import "fmt"

// ErrAdminRequired is returned when administrator privileges are required
type ErrAdminRequired struct {
	Message string
}

func (e *ErrAdminRequired) Error() string {
	return fmt.Sprintf("administrator privileges required: %s", e.Message)
}

// NewErrAdminRequired creates a new ErrAdminRequired error
func NewErrAdminRequired(message string) error {
	return &ErrAdminRequired{
		Message: message,
	}
}

// ErrCleanupFailed is returned when a cleanup operation fails
type ErrCleanupFailed struct {
	CleanerName string
	Message     string
}

func (e *ErrCleanupFailed) Error() string {
	return fmt.Sprintf("cleanup failed for %s: %s", e.CleanerName, e.Message)
}

// NewErrCleanupFailed creates a new ErrCleanupFailed error
func NewErrCleanupFailed(cleanerName, message string) error {
	return &ErrCleanupFailed{
		CleanerName: cleanerName,
		Message:     message,
	}
}

// ErrNotSupported is returned when an operation is not supported
type ErrNotSupported struct {
	Operation string
	Reason    string
}

func (e *ErrNotSupported) Error() string {
	return fmt.Sprintf("operation not supported: %s - %s", e.Operation, e.Reason)
}

// NewErrNotSupported creates a new ErrNotSupported error
func NewErrNotSupported(operation, reason string) error {
	return &ErrNotSupported{
		Operation: operation,
		Reason:    reason,
	}
}
