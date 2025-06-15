package reporter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CacheReporter handles reporting of cache sizes
type CacheReporter struct{}

// NewCacheReporter creates a new CacheReporter instance
func NewCacheReporter() *CacheReporter {
	return &CacheReporter{}
}

// getDirSize calculates the size of a directory in bytes
func (r *CacheReporter) getDirSize(path string) (int64, error) {
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

// formatSize converts bytes to human-readable format
func (r *CacheReporter) formatSize(size int64) string {
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

// GetNPMCacheSize returns the size of npm cache
func (r *CacheReporter) GetNPMCacheSize() (string, error) {
	npmCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "npm-cache")
	if _, err := os.Stat(npmCache); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := r.getDirSize(npmCache)
	if err != nil {
		return "Error", err
	}
	return r.formatSize(size), nil
}

// GetYarnCacheSize returns the size of yarn cache
func (r *CacheReporter) GetYarnCacheSize() (string, error) {
	yarnCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "Yarn", "cache", "v6")
	if _, err := os.Stat(yarnCache); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := r.getDirSize(yarnCache)
	if err != nil {
		return "Error", err
	}
	return r.formatSize(size), nil
}

// GetDockerCacheSize returns the size of Docker cache
func (r *CacheReporter) GetDockerCacheSize() (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "Not installed", nil
	}

	// Check if Docker daemon is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return "Docker not running", nil
	}

	cmd = exec.Command("docker", "system", "df", "--format", "{{.Size}}")
	output, err := cmd.Output()
	if err != nil {
		return "Error getting size", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// GetWinSxSTempSize returns the size of WinSxS temp folder
func (r *CacheReporter) GetWinSxSTempSize() (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	winsxsTemp := filepath.Join(os.Getenv("WINDIR"), "WinSxS", "Temp")
	if _, err := os.Stat(winsxsTemp); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := r.getDirSize(winsxsTemp)
	if err != nil {
		return "Error", err
	}
	return r.formatSize(size), nil
}

// GetWindowsTempSize returns the size of Windows temporary files
func (r *CacheReporter) GetWindowsTempSize() (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	tempDir := os.TempDir()
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := r.getDirSize(tempDir)
	if err != nil {
		return "Error", err
	}
	return r.formatSize(size), nil
}

// GetWindowsChunkSize returns the size of Windows chunk files
func (r *CacheReporter) GetWindowsChunkSize() (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	chunkDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "WER", "ReportQueue")
	if _, err := os.Stat(chunkDir); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := r.getDirSize(chunkDir)
	if err != nil {
		return "Error", err
	}
	return r.formatSize(size), nil
}

// GetCacheSizes returns a map of all cache sizes
func (r *CacheReporter) GetCacheSizes() map[string]string {
	sizes := make(map[string]string)

	// NPM Cache
	npmSize, err := r.GetNPMCacheSize()
	if err != nil {
		sizes["npm cache"] = "Error"
	} else {
		sizes["npm cache"] = npmSize
	}

	// Yarn Cache
	yarnSize, err := r.GetYarnCacheSize()
	if err != nil {
		sizes["yarn cache"] = "Error"
	} else {
		sizes["yarn cache"] = yarnSize
	}

	// Docker Cache
	dockerSize, err := r.GetDockerCacheSize()
	if err != nil {
		sizes["docker cache"] = "Error"
	} else {
		sizes["docker cache"] = dockerSize
	}

	// WinSxS Temp
	winsxsSize, err := r.GetWinSxSTempSize()
	if err != nil {
		sizes["WinSxS temp"] = "Error"
	} else {
		sizes["WinSxS temp"] = winsxsSize
	}

	// Windows Temp
	winTempSize, err := r.GetWindowsTempSize()
	if err != nil {
		sizes["Windows temp"] = "Error"
	} else {
		sizes["Windows temp"] = winTempSize
	}

	// Windows Chunks
	winChunkSize, err := r.GetWindowsChunkSize()
	if err != nil {
		sizes["Windows chunks"] = "Error"
	} else {
		sizes["Windows chunks"] = winChunkSize
	}

	return sizes
} 