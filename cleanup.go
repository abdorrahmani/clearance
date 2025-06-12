package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/gookit/color"
)

// checkAdminPrivileges checks if the program is running with administrator privileges
func checkAdminPrivileges() error {
	if runtime.GOOS == "windows" {
		// Try to open a system file that requires admin rights
		_, err := syscall.Open("\\\\.\\PHYSICALDRIVE0", syscall.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("this program requires administrator privileges")
		}
	}
	return nil
}

// cleanNPMCache attempts to clean npm cache using multiple methods
func cleanNPMCache() error {
	fmt.Println("[npm] Attempting to remove npm cache folder...")
	npmCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "npm-cache")
	if err := os.RemoveAll(npmCache); err == nil {
		fmt.Println("[npm] Folder removed successfully.")
		return nil
	} else {
		fmt.Printf("[npm] Folder removal failed: %v\n", err)
	}

	fmt.Println("[npm] Fallback: running 'npm cache clean --force'...")
	if npmPath, err := exec.LookPath("npm"); err == nil {
		cmd := exec.Command(npmPath, "cache", "clean", "--force")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("[npm] npm CLI cache clean succeeded.")
			return nil
		} else {
			fmt.Printf("[npm] npm CLI cache clean failed: %v\n", err)
		}
	} else {
		fmt.Println("[npm] npm not found in PATH.")
	}

	return fmt.Errorf("failed to clean npm cache using both direct deletion and npm CLI")
}

// cleanYarnCache attempts to clean yarn cache using multiple methods
func cleanYarnCache() error {
	fmt.Println("[yarn] Attempting to remove yarn cache folder...")
	yarnCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "Yarn", "cache", "v6")
	if err := os.RemoveAll(yarnCache); err == nil {
		fmt.Println("[yarn] Folder removed successfully.")
		return nil
	} else {
		fmt.Printf("[yarn] Folder removal failed: %v\n", err)
	}

	fmt.Println("[yarn] Fallback: running 'yarn cache clean'...")
	if yarnPath, err := exec.LookPath("yarn"); err == nil {
		cmd := exec.Command(yarnPath, "cache", "clean")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("[yarn] yarn CLI cache clean succeeded.")
			return nil
		} else {
			fmt.Printf("[yarn] yarn CLI cache clean failed: %v\n", err)
		}
	} else {
		fmt.Println("[yarn] yarn not found in PATH.")
	}

	return fmt.Errorf("failed to clean yarn cache using both direct deletion and yarn CLI")
}

// cleanDockerCache executes Docker cleanup commands
func cleanDockerCache() error {
	fmt.Println("[docker] Running Docker cleanup commands...")
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("[docker] Docker not found in PATH.")
		return fmt.Errorf("docker not found in PATH")
	}

	cmd := exec.Command("docker", "system", "prune", "--all", "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[docker] Failed to run docker system prune: %v\n", err)
		return fmt.Errorf("failed to run docker system prune: %v", err)
	}

	cmd = exec.Command("docker", "volume", "prune", "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[docker] Failed to run docker volume prune: %v\n", err)
		return fmt.Errorf("failed to run docker volume prune: %v", err)
	}

	cmd = exec.Command("docker", "builder", "prune", "--all", "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[docker] Failed to run docker builder prune: %v\n", err)
		return fmt.Errorf("failed to run docker builder prune: %v", err)
	}

	fmt.Println("[docker] Docker cleanup completed.")
	return nil
}

// cleanWinSxSTemp cleans Windows WinSxS temp files
func cleanWinSxSTemp() error {
	fmt.Println("[winsxs] Attempting to clean WinSxS Temp folder...")
	if runtime.GOOS != "windows" {
		fmt.Println("[winsxs] WinSxS cleanup is only available on Windows.")
		return fmt.Errorf("WinSxS cleanup is only available on Windows")
	}

	winsxsTemp := filepath.Join(os.Getenv("WINDIR"), "WinSxS", "Temp")

	// First try: PowerShell command with elevated privileges and better error handling
	fmt.Println("[winsxs] Attempting to clean using PowerShell...")
	psCmd := exec.Command("powershell", "-Command", `
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

	// Second try: Manual cleanup with better error handling
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

// getDirSize calculates the size of a directory in bytes
func getDirSize(path string) (int64, error) {
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
func formatSize(size int64) string {
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

// getNPMCacheSize returns the size of npm cache
func getNPMCacheSize() (string, error) {
	npmCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "npm-cache")
	if _, err := os.Stat(npmCache); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := getDirSize(npmCache)
	if err != nil {
		return "Error", err
	}
	return formatSize(size), nil
}

// getYarnCacheSize returns the size of yarn cache
func getYarnCacheSize() (string, error) {
	yarnCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "Yarn", "Cache")
	if _, err := os.Stat(yarnCache); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := getDirSize(yarnCache)
	if err != nil {
		return "Error", err
	}
	return formatSize(size), nil
}

// getDockerCacheSize returns the size of Docker cache
func getDockerCacheSize() (string, error) {
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

// getWinSxSTempSize returns the size of WinSxS temp folder
func getWinSxSTempSize() (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	winsxsTemp := filepath.Join(os.Getenv("WINDIR"), "WinSxS", "Temp")
	if _, err := os.Stat(winsxsTemp); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := getDirSize(winsxsTemp)
	if err != nil {
		return "Error", err
	}
	return formatSize(size), nil
}

// getWindowsTempSize returns the size of Windows temporary files
func getWindowsTempSize() (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	tempDir := os.TempDir()
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := getDirSize(tempDir)
	if err != nil {
		return "Error", err
	}
	return formatSize(size), nil
}

// getWindowsChunkSize returns the size of Windows chunk files
func getWindowsChunkSize() (string, error) {
	if runtime.GOOS != "windows" {
		return "N/A", nil
	}

	chunkDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "WER", "ReportQueue")
	if _, err := os.Stat(chunkDir); os.IsNotExist(err) {
		return "Not found", nil
	}
	size, err := getDirSize(chunkDir)
	if err != nil {
		return "Error", err
	}
	return formatSize(size), nil
}

// cleanupWindowsTemp cleans Windows temporary files
func cleanupWindowsTemp() error {
	fmt.Println("[wintemp] Attempting to clean Windows temporary files...")
	if runtime.GOOS != "windows" {
		fmt.Println("[wintemp] Windows temp cleanup is only available on Windows.")
		return fmt.Errorf("windows temp cleanup is only available on Windows")
	}

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

// cleanupWindowsChunks cleans Windows error reporting chunks
func cleanupWindowsChunks() error {
	fmt.Println("[winchunks] Attempting to clean Windows error reporting chunks...")
	if runtime.GOOS != "windows" {
		fmt.Println("[winchunks] Windows chunks cleanup is only available on Windows.")
		return fmt.Errorf("windows chunks cleanup is only available on Windows")
	}

	chunkDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "WER", "ReportQueue")
	if _, err := os.Stat(chunkDir); os.IsNotExist(err) {
		fmt.Println("[winchunks] ReportQueue directory not found.")
		return nil
	}

	// First try: PowerShell command with elevated privileges
	fmt.Println("[winchunks] Attempting to clean using PowerShell...")
	psCmd := exec.Command("powershell", "-Command", `
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

// reportCacheSizes displays the size of all caches
func reportCacheSizes() error {
	color.Blue.Println("\n📊 Cache Size Report")
	color.Blue.Println("===================")

	// NPM Cache
	npmSize, err := getNPMCacheSize()
	if err != nil {
		color.Red.Printf("npm cache: Error - %v\n", err)
	} else {
		if npmSize == "Not found" {
			color.Yellow.Printf("npm cache: %s\n", npmSize)
		} else {
			color.Green.Printf("npm cache: %s\n", npmSize)
		}
	}

	// Yarn Cache
	yarnSize, err := getYarnCacheSize()
	if err != nil {
		color.Red.Printf("yarn cache: Error - %v\n", err)
	} else {
		if yarnSize == "Not found" {
			color.Yellow.Printf("yarn cache: %s\n", yarnSize)
		} else {
			color.Green.Printf("yarn cache: %s\n", yarnSize)
		}
	}

	// Docker Cache
	dockerSize, err := getDockerCacheSize()
	if err != nil {
		color.Red.Printf("docker cache: Error - %v\n", err)
	} else {
		if dockerSize == "Not installed" || dockerSize == "Docker not running" {
			color.Yellow.Printf("docker cache: %s\n", dockerSize)
		} else {
			color.Green.Printf("docker cache: %s\n", dockerSize)
		}
	}

	// WinSxS Temp
	winsxsSize, err := getWinSxSTempSize()
	if err != nil {
		color.Red.Printf("WinSxS temp: Error - %v\n", err)
	} else {
		if winsxsSize == "Not found" || winsxsSize == "N/A" {
			color.Yellow.Printf("WinSxS temp: %s\n", winsxsSize)
		} else {
			color.Green.Printf("WinSxS temp: %s\n", winsxsSize)
		}
	}

	// Windows Temp
	winTempSize, err := getWindowsTempSize()
	if err != nil {
		color.Red.Printf("Windows temp: Error - %v\n", err)
	} else {
		if winTempSize == "Not found" || winTempSize == "N/A" {
			color.Yellow.Printf("Windows temp: %s\n", winTempSize)
		} else {
			color.Green.Printf("Windows temp: %s\n", winTempSize)
		}
	}

	// Windows Chunks
	winChunkSize, err := getWindowsChunkSize()
	if err != nil {
		color.Red.Printf("Windows chunks: Error - %v\n", err)
	} else {
		if winChunkSize == "Not found" || winChunkSize == "N/A" {
			color.Yellow.Printf("Windows chunks: %s\n", winChunkSize)
		} else {
			color.Green.Printf("Windows chunks: %s\n", winChunkSize)
		}
	}

	return nil
}
