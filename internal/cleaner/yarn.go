package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// YarnCleaner handles cleaning of yarn cache
type YarnCleaner struct {
	*BaseCleaner
}

// NewYarnCleaner creates a new YarnCleaner
func NewYarnCleaner() *YarnCleaner {
	return &YarnCleaner{
		BaseCleaner: NewBaseCleaner("yarn"),
	}
}

// Clean performs the yarn cache cleaning operation
func (y *YarnCleaner) Clean(ctx context.Context) error {
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
		cmd := exec.CommandContext(ctx, yarnPath, "cache", "clean")
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

// GetSize returns the size of yarn cache
func (y *YarnCleaner) GetSize(ctx context.Context) (string, error) {
	yarnCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "Yarn", "Cache")
	exists, err := CheckPathExists(yarnCache)
	if err != nil {
		return "Error", err
	}
	if !exists {
		return "Not found", nil
	}
	size, err := GetDirSize(yarnCache)
	if err != nil {
		return "Error", err
	}
	return FormatSize(size), nil
}
