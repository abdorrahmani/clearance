package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// WindowsCleaner handles cleaning of Windows system files
type WindowsCleaner struct {
	*BaseCleaner
	cleanType string
}

// NewWindowsCleaner creates a new WindowsCleaner
func NewWindowsCleaner(cleanType string) *WindowsCleaner {
	return &WindowsCleaner{
		BaseCleaner: NewBaseCleaner(cleanType),
		cleanType:   cleanType,
	}
}

// Clean performs the Windows system cleaning operation
func (w *WindowsCleaner) Clean(ctx context.Context) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("%s cleanup is only available on Windows", w.cleanType)
	}

	switch w.cleanType {
	case "winsxs":
		return w.cleanWinSxS(ctx)
	case "wintemp":
		return w.cleanWindowsTemp(ctx)
	case "winchunks":
		return w.cleanWindowsChunks(ctx)
	default:
		return fmt.Errorf("unknown Windows cleaner type: %s", w.cleanType)
	}
}

// GetSize returns the size of Windows system files
func (w *WindowsCleaner) GetSize(ctx context.Context) (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	var path string
	switch w.cleanType {
	case "winsxs":
		path = filepath.Join(os.Getenv("WINDIR"), "WinSxS", "Temp")
	case "wintemp":
		path = os.TempDir()
	case "winchunks":
		path = filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "WER", "ReportQueue")
	default:
		return "N/A", fmt.Errorf("unknown Windows cleaner type: %s", w.cleanType)
	}

	exists, err := CheckPathExists(path)
	if err != nil {
		return "Error", err
	}
	if !exists {
		return "Not found", nil
	}
	size, err := GetDirSize(path)
	if err != nil {
		return "Error", err
	}
	return FormatSize(size), nil
}

func (w *WindowsCleaner) cleanWinSxS(ctx context.Context) error {
	fmt.Println("[winsxs] Attempting to clean WinSxS Temp folder...")
	winsxsTemp := filepath.Join(os.Getenv("WINDIR"), "WinSxS", "Temp")

	// First try: PowerShell command with elevated privileges
	fmt.Println("[winsxs] Attempting to clean using PowerShell...")
	psCmd := exec.CommandContext(ctx, "powershell", "-Command", `
		$ErrorActionPreference = 'Stop'
		$paths = @(
			'InFlight',
			'PendingDeletes',
			'PendingRenames'
		)
		foreach ($path in $paths) {
			$fullPath = Join-Path $env:WINDIR "WinSxS\Temp\$path"
			if (Test-Path $fullPath) {
				try {
					Get-ChildItem -Path $fullPath -Recurse | Remove-Item -Force -Recurse -ErrorAction Stop
					Write-Host "[winsxs] Successfully cleaned $path"
				} catch {
					Write-Host "[winsxs] Warning: Could not clean $path - $($_.Exception.Message)"
				}
			}
		}
	`)
	psCmd.Stdout = os.Stdout
	psCmd.Stderr = os.Stderr
	if err := psCmd.Run(); err != nil {
		fmt.Printf("[winsxs] PowerShell cleanup encountered issues: %v\n", err)
	}

	// Second try: Manual cleanup
	fmt.Println("[winsxs] Attempting manual cleanup...")
	entries, err := os.ReadDir(winsxsTemp)
	if err != nil {
		return fmt.Errorf("failed to read WinSxS Temp folder: %w", err)
	}

	var total, failed int
	for _, entry := range entries {
		total++
		path := filepath.Join(winsxsTemp, entry.Name())

		// Skip system-protected folders
		if entry.Name() == "InFlight" || entry.Name() == "PendingDeletes" || entry.Name() == "PendingRenames" {
			fmt.Printf("[winsxs] Skipping system-protected folder: %s\n", entry.Name())
			continue
		}

		err := os.RemoveAll(path)
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("[winsxs] Access denied for: %s (This is normal for system-protected files)\n", path)
			} else {
				fmt.Printf("[winsxs] Failed to remove: %s (%v)\n", path, err)
			}
			failed++
		} else {
			fmt.Printf("[winsxs] Successfully removed: %s\n", path)
		}
	}

	if failed == 0 {
		fmt.Println("[winsxs] WinSxS Temp folder cleaned successfully.")
		return nil
	}

	if failed < total {
		fmt.Printf("[winsxs] Partial cleanup complete. %d of %d items removed successfully.\n", total-failed, total)
		fmt.Println("[winsxs] Some files could not be deleted due to system protection. This is normal for active Windows Update operations.")
		return nil
	}

	return fmt.Errorf("failed to clean any files in WinSxS Temp folder")
}

func (w *WindowsCleaner) cleanWindowsTemp(ctx context.Context) error {
	fmt.Println("[wintemp] Attempting to clean Windows temporary files...")
	tempDir := os.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	var total, failed int
	for _, entry := range entries {
		total++
		path := filepath.Join(tempDir, entry.Name())

		// Skip if the file is currently in use
		if file, err := os.OpenFile(path, os.O_RDWR, 0); err == nil {
			file.Close()
			err := os.RemoveAll(path)
			if err != nil {
				if os.IsPermission(err) {
					fmt.Printf("[wintemp] Access denied for: %s (File might be in use)\n", path)
				} else {
					fmt.Printf("[wintemp] Failed to remove: %s (%v)\n", path, err)
				}
				failed++
			} else {
				fmt.Printf("[wintemp] Successfully removed: %s\n", path)
			}
		} else {
			fmt.Printf("[wintemp] Skipping in-use file: %s\n", path)
			failed++
		}
	}

	if failed == 0 {
		fmt.Println("[wintemp] Windows temp folder cleaned successfully.")
		return nil
	}

	if failed < total {
		fmt.Printf("[wintemp] Partial cleanup complete. %d of %d items removed successfully.\n", total-failed, total)
		fmt.Println("[wintemp] Some files could not be deleted as they are currently in use.")
		return nil
	}

	return fmt.Errorf("failed to clean any files in Windows temp folder")
}

func (w *WindowsCleaner) cleanWindowsChunks(ctx context.Context) error {
	fmt.Println("[winchunks] Attempting to clean Windows error reporting chunks...")
	chunkDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "WER", "ReportQueue")
	if _, err := os.Stat(chunkDir); os.IsNotExist(err) {
		fmt.Println("[winchunks] ReportQueue directory not found.")
		return nil
	}

	// First try: PowerShell command with elevated privileges
	fmt.Println("[winchunks] Attempting to clean using PowerShell...")
	psCmd := exec.CommandContext(ctx, "powershell", "-Command", `
		$ErrorActionPreference = 'Stop'
		$chunkDir = Join-Path $env:LOCALAPPDATA "Microsoft\Windows\WER\ReportQueue"
		if (Test-Path $chunkDir) {
			try {
				Get-ChildItem -Path $chunkDir -Recurse | Remove-Item -Force -Recurse -ErrorAction Stop
				Write-Host "[winchunks] Successfully cleaned ReportQueue"
			} catch {
				Write-Host "[winchunks] Warning: Could not clean ReportQueue - $($_.Exception.Message)"
			}
		}
	`)
	psCmd.Stdout = os.Stdout
	psCmd.Stderr = os.Stderr
	if err := psCmd.Run(); err != nil {
		fmt.Printf("[winchunks] PowerShell cleanup encountered issues: %v\n", err)
	}

	// Second try: Manual cleanup
	fmt.Println("[winchunks] Attempting manual cleanup...")
	entries, err := os.ReadDir(chunkDir)
	if err != nil {
		return fmt.Errorf("failed to read ReportQueue directory: %w", err)
	}

	var total, failed int
	for _, entry := range entries {
		total++
		path := filepath.Join(chunkDir, entry.Name())

		err := os.RemoveAll(path)
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("[winchunks] Access denied for: %s (This is normal for system-protected files)\n", path)
			} else {
				fmt.Printf("[winchunks] Failed to remove: %s (%v)\n", path, err)
			}
			failed++
		} else {
			fmt.Printf("[winchunks] Successfully removed: %s\n", path)
		}
	}

	if failed == 0 {
		fmt.Println("[winchunks] Windows chunks cleaned successfully.")
		return nil
	}

	if failed < total {
		fmt.Printf("[winchunks] Partial cleanup complete. %d of %d items removed successfully.\n", total-failed, total)
		fmt.Println("[winchunks] Some files could not be deleted due to system protection.")
		return nil
	}

	return fmt.Errorf("failed to clean any files in Windows chunks folder")
}
