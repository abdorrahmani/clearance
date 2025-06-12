package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

const clearanceLogo = `
  C L E A R A N C E
  =================
`

var (
	cleanNPM    bool
	cleanYarn   bool
	cleanDocker bool
	cleanWinSxS bool
	cleanAll    bool
	reportSize  bool
	version     = "0.1.0"
	commit      = "none"
	date        = "2025-06-09"
)

var rootCmd = &cobra.Command{
	Use:   "clearance",
	Short: "A lightweight CLI tool to clean up development caches",
	Long: `Clearance is a CLI tool that helps free up disk space by cleaning various development caches.
It can clean npm, yarn, Docker, and Windows system temp files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cleanAll {
			cleanNPM = true
			cleanYarn = true
			cleanDocker = true
			cleanWinSxS = true
		}

		if !cleanNPM && !cleanYarn && !cleanDocker && !cleanWinSxS {
			return fmt.Errorf("no cleanup options specified. Use --help for available options")
		}

		if err := checkAdminPrivileges(); err != nil {
			fmt.Printf("[error] %v\n", err)
			return err
		}

		var errors []error

		if cleanNPM {
			if err := cleanNPMCache(); err != nil {
				fmt.Printf("[error] %v\n", err)
				errors = append(errors, err)
			}
		}

		if cleanYarn {
			if err := cleanYarnCache(); err != nil {
				fmt.Printf("[error] %v\n", err)
				errors = append(errors, err)
			}
		}

		if cleanDocker {
			if err := cleanDockerCache(); err != nil {
				fmt.Printf("[error] %v\n", err)
				errors = append(errors, err)
			}
		}

		if cleanWinSxS {
			if err := cleanWinSxSTemp(); err != nil {
				fmt.Printf("[error] %v\n", err)
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			fmt.Printf("\n⚠️  Clearance completed with %d error(s). Some operations may have failed.\n", len(errors))
			return fmt.Errorf("some cleanup operations failed")
		}

		fmt.Println("\n✅ Clearance finished successfully!")
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVar(&cleanNPM, "npm", false, "Clean npm cache")
	rootCmd.Flags().BoolVar(&cleanYarn, "yarn", false, "Clean yarn cache")
	rootCmd.Flags().BoolVar(&cleanDocker, "docker", false, "Clean Docker cache")
	rootCmd.Flags().BoolVar(&cleanWinSxS, "winsxs", false, "Clean WinSxS temp files")
	rootCmd.Flags().BoolVar(&cleanAll, "all", false, "Clean all caches")
	rootCmd.Flags().BoolVar(&reportSize, "report", false, "Report cache sizes")
}

func showVersion() {
	fmt.Printf("\nClearance v%s\n", version)
	fmt.Printf("Build: %s (%s)\n", commit, date)
	fmt.Printf("OS/Arch: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
}

func showAdminWarning() {
	color.Red.Print("\n🔒 Administrator Privileges Required 🔒\n")
	color.Red.Println("========================================")
	color.Yellow.Println("⚠️  This tool requires administrator privileges to clean system caches.")
	color.Yellow.Println("Please run this tool as administrator.")
	color.Yellow.Println("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

func showMenu() {
	color.Blue.Print(clearanceLogo)
	color.Blue.Println("         Clearance - Cache Cleanup Tool")
	color.Blue.Println("=========================================")
	color.Bold.Println("\n📁 Available Options:")
	color.Green.Println("1. Clean npm cache")
	color.Green.Println("2. Clean yarn cache")
	color.Green.Println("3. Clean Docker cache")
	color.Green.Println("4. Clean WinSxS temp files")
	color.Green.Println("5. Clean everything")
	color.Green.Println("6. Report cache sizes")
	color.Green.Println("7. Exit")
	color.Green.Println("8. Show version")
	color.Bold.Print("\n→ Please enter your choice (1-8): ")
}

func runCleanup(option string) error {
	var errors []error

	switch option {
	case "1":
		if err := cleanNPMCache(); err != nil {
			errors = append(errors, err)
		}
	case "2":
		if err := cleanYarnCache(); err != nil {
			errors = append(errors, err)
		}
	case "3":
		if err := cleanDockerCache(); err != nil {
			errors = append(errors, err)
		}
	case "4":
		if err := cleanWinSxSTemp(); err != nil {
			errors = append(errors, err)
		}
	case "5":
		if err := cleanNPMCache(); err != nil {
			errors = append(errors, err)
		}
		if err := cleanYarnCache(); err != nil {
			errors = append(errors, err)
		}
		if err := cleanDockerCache(); err != nil {
			errors = append(errors, err)
		}
		if err := cleanWinSxSTemp(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		color.Yellow.Printf("\n⚠️  Cleanup completed with %d error(s). Some operations may have failed.\n", len(errors))
		return fmt.Errorf("some cleanup operations failed")
	}

	color.Green.Println("\n✅ Cleanup finished successfully!")
	return nil
}

func main() {
	// Check for admin privileges first
	if err := checkAdminPrivileges(); err != nil {
		showAdminWarning()
	}

	// If running from PowerShell, set up the environment
	if runtime.GOOS == "windows" {
		// Set VirtualTerminalLevel for ANSI color support
		cmd := exec.Command("powershell", "-NoExit", "-Command",
			"Set-ItemProperty -Path 'HKCU:\\Console' -Name 'VirtualTerminalLevel' -Value 1; "+
				"$host.UI.RawUI.WindowTitle = 'Clearance - Cache Cleanup Tool'")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Printf("Warning: Could not set up PowerShell environment: %v\n", err)
		}
	}

	for {
		showMenu()
		reader := bufio.NewReader(os.Stdin)
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1", "2", "3", "4", "5":
			if err := runCleanup(option); err != nil {
				color.Red.Printf("[error] %v\n", err)
			}
			color.Yellow.Println("\nPress Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		case "6":
			if err := reportCacheSizes(); err != nil {
				color.Red.Printf("[error] %v\n", err)
			}
			color.Yellow.Println("\nPress Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		case "7":
			color.Blue.Println("\n👋 Thank you for using Clearance!")
			os.Exit(0)
		case "8":
			showVersion()
			color.Yellow.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		default:
			color.Red.Println("\n❌ Invalid option. Please try again.")
			color.Yellow.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}
