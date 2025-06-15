package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/abdorrahmani/clearance/internal/cleaner"
	"github.com/abdorrahmani/clearance/internal/reporter"
	"github.com/abdorrahmani/clearance/internal/ui"
	"github.com/abdorrahmani/clearance/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "none"
	date    = "2025-06-09"
)

func executeCleanup(ui *ui.UI, options []string) error {
	ctx := context.Background()

	ui.ShowSelectedOptions(options)

	// Handle cache size reporting
	if len(options) == 1 && (options[0] == "7" || options[0] == "report") {
		reporter := reporter.NewCacheReporter()
		sizes := reporter.GetCacheSizes()
		ui.ShowCacheSizeReport(sizes)
		return nil
	}

	cleaners := []cleaner.Cleaner{}
	for _, opt := range options {
		switch opt {
		case "1", "npm":
			cleaners = append(cleaners, cleaner.NewNPMCleaner())
		case "2", "yarn":
			cleaners = append(cleaners, cleaner.NewYarnCleaner())
		case "3", "docker":
			cleaners = append(cleaners, cleaner.NewDockerCleaner())
		case "4", "winsxs":
			cleaners = append(cleaners, cleaner.NewWindowsCleaner("winsxs"))
		case "5", "wintemp":
			cleaners = append(cleaners, cleaner.NewWindowsCleaner("wintemp"))
		case "6", "winchunks":
			cleaners = append(cleaners, cleaner.NewWindowsCleaner("winchunks"))
		case "7", "report":
			reporter := reporter.NewCacheReporter()
			sizes := reporter.GetCacheSizes()
			ui.ShowCacheSizeReport(sizes)
			return nil
		case "8", "exit":
			return nil
		}
	}

	if len(cleaners) == 0 {
		return errors.NewErrNotSupported("cleanup", "no valid cleanup options selected")
	}

	if err := cleaner.CheckAdminPrivileges(); err != nil {
		ui.ShowAdminWarning()
		return err
	}

	ui.ShowCleanupStart()

	var errs []error
	for _, c := range cleaners {
		if err := c.Clean(ctx); err != nil {
			ui.ShowError(err)
			errs = append(errs, err)
		}
	}

	ui.ShowCleanupComplete(len(errs))
	if len(errs) > 0 {
		return errors.NewErrCleanupFailed("all", "some cleanup operations failed")
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "clearance",
	Short: "A lightweight CLI tool to clean up development caches",
	Long: `Clearance is a CLI tool that helps free up disk space by cleaning various development caches.
It can clean npm, yarn, Docker, and Windows system temp files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui := ui.NewUI()

		for {
			ui.ShowMenu()
			input := ui.ReadInput()
			if input == "" {
				continue
			}

			var options []string
			if strings.ToLower(input) == "all" {
				options = []string{"1", "2", "3", "4", "5", "6"}
			} else if strings.ToLower(input) == "exit" {
				options = []string{"8"}
			} else {
				options = strings.Split(input, ",")
			}

			if err := executeCleanup(ui, options); err != nil {
				ui.ShowError(err)
			}

			ui.WaitForEnter()
		}
	},
}

func main() {
	// If running from PowerShell, set up the environment
	if runtime.GOOS == "windows" {
		// Set VirtualTerminalLevel for ANSI color support
		cmd := exec.Command("powershell", "-Command",
			"Set-ItemProperty -Path 'HKCU:\\Console' -Name 'VirtualTerminalLevel' -Value 1; "+
				"$host.UI.RawUI.WindowTitle = 'Clearance - Cache Cleanup Tool'")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			ui := ui.NewUI()
			ui.ShowWarning(fmt.Sprintf("Could not set up PowerShell environment: %v", err))
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
