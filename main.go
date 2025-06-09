package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cleanNPM    bool
	cleanYarn   bool
	cleanDocker bool
	cleanWinSxS bool
	cleanAll    bool
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
}

func showAdminWarning() {
	fmt.Println("⚠️  This tool requires administrator privileges to clean system caches.")
	fmt.Println("Please run this tool as administrator.")
	fmt.Println("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

func showMenu() {
	fmt.Println("\n🧹 Clearance - Cache Cleanup Tool")
	fmt.Println("================================")
	fmt.Println("1. Clean npm cache")
	fmt.Println("2. Clean yarn cache")
	fmt.Println("3. Clean Docker cache")
	fmt.Println("4. Clean WinSxS temp files")
	fmt.Println("5. Clean everything")
	fmt.Println("6. Exit")
	fmt.Print("\nSelect an option (1-6): ")
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
		fmt.Printf("\n⚠️  Cleanup completed with %d error(s). Some operations may have failed.\n", len(errors))
		return fmt.Errorf("some cleanup operations failed")
	}

	fmt.Println("\n✅ Cleanup finished successfully!")
	return nil
}

func main() {
	// Check for admin privileges first
	if err := checkAdminPrivileges(); err != nil {
		showAdminWarning()
	}

	// If running from PowerShell, keep the window open
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoExit", "-Command", "& { $host.UI.RawUI.WindowTitle = 'Clearance - Cache Cleanup Tool' }")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}

	for {
		showMenu()
		reader := bufio.NewReader(os.Stdin)
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1", "2", "3", "4", "5":
			if err := runCleanup(option); err != nil {
				fmt.Printf("[error] %v\n", err)
			}
			fmt.Println("\nPress Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		case "6":
			fmt.Println("\n👋 Thank you for using Clearance!")
			os.Exit(0)
		default:
			fmt.Println("\n❌ Invalid option. Please try again.")
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}
