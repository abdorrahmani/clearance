package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
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

	// First try: PowerShell command with elevated privileges
	fmt.Println("[winsxs] Attempting to clean using PowerShell...")
	psCmd := exec.Command("powershell", "-Command", fmt.Sprintf("Remove-Item -Recurse -Force '%s\\*' -ErrorAction SilentlyContinue", winsxsTemp))
	psCmd.Stdout = os.Stdout
	psCmd.Stderr = os.Stderr
	if err := psCmd.Run(); err == nil {
		fmt.Println("[winsxs] PowerShell cleanup completed.")
		return nil
	}

	// Second try: Manual cleanup of individual items
	fmt.Println("[winsxs] PowerShell cleanup failed, attempting manual cleanup...")
	entries, err := os.ReadDir(winsxsTemp)
	if err != nil {
		return fmt.Errorf("failed to read WinSxS Temp folder: %w", err)
	}

	var total, failed int
	for _, entry := range entries {
		total++
		path := filepath.Join(winsxsTemp, entry.Name())
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Printf("[winsxs] Failed to remove: %s (%v)\n", path, err)
			failed++
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
