package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// NPMCleaner handles cleaning of npm cache
type NPMCleaner struct {
	*BaseCleaner
}

// NewNPMCleaner creates a new NPMCleaner
func NewNPMCleaner() *NPMCleaner {
	return &NPMCleaner{
		BaseCleaner: NewBaseCleaner("npm"),
	}
}

// Clean performs the npm cache cleaning operation
func (n *NPMCleaner) Clean(ctx context.Context) error {
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
		cmd := exec.CommandContext(ctx, npmPath, "cache", "clean", "--force")
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

// GetSize returns the size of npm cache
func (n *NPMCleaner) GetSize(ctx context.Context) (string, error) {
	npmCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "npm-cache")
	exists, err := CheckPathExists(npmCache)
	if err != nil {
		return "Error", err
	}
	if !exists {
		return "Not found", nil
	}
	size, err := GetDirSize(npmCache)
	if err != nil {
		return "Error", err
	}
	return FormatSize(size), nil
}
