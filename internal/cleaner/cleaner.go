package cleaner

import (
	"context"
	"os"
	"runtime"

	"github.com/abdorrahmani/clearance/pkg/errors"
)

// Cleaner defines the interface for cache cleaning operations
type Cleaner interface {
	// Clean performs the cache cleaning operation
	Clean(ctx context.Context) error
	// GetSize returns the current size of the cache
	GetSize(ctx context.Context) (string, error)
	// GetName returns the name of the cleaner
	GetName() string
}

// CleanResult represents the result of a cleaning operation
type CleanResult struct {
	CleanerName string
	Error       error
	SizeBefore  string
	SizeAfter   string
}

// CleanOptions represents the options for cleaning operations
type CleanOptions struct {
	CleanNPM           bool
	CleanYarn          bool
	CleanDocker        bool
	CleanWinSxS        bool
	CleanWindowsTemp   bool
	CleanWindowsChunks bool
	CleanAll           bool
	ReportSize         bool
}

// CheckAdminPrivileges checks if the program is running with administrator privileges
func CheckAdminPrivileges() error {
	if runtime.GOOS == "windows" {
		// Try to open a system file that requires admin rights
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			return errors.NewErrAdminRequired("this program requires administrator privileges")
		}
	}
	return nil
}
