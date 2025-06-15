package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
)

// BaseCleaner provides common functionality for all cleaners
type BaseCleaner struct {
	name string
}

// NewBaseCleaner creates a new BaseCleaner
func NewBaseCleaner(name string) *BaseCleaner {
	return &BaseCleaner{
		name: name,
	}
}

// GetName returns the name of the cleaner
func (b *BaseCleaner) GetName() string {
	return b.name
}

// GetDirSize calculates the size of a directory in bytes
func GetDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// FormatSize converts bytes to human-readable format
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// CheckPathExists checks if a path exists
func CheckPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
