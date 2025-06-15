package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

const clearanceLogo = `
  C L E A R A N C E
  =================
`

// UI handles all user interface interactions
type UI struct {
	reader *bufio.Reader
}

// NewUI creates a new UI instance
func NewUI() *UI {
	return &UI{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ClearScreen clears the terminal screen
func (u *UI) ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

// ShowHeader displays the application header
func (u *UI) ShowHeader(version string) {
	color.Blue.Println("\n╔════════════════════════════════════════════════════════════╗")
	color.Blue.Println("║                    Clearance v" + version + "                    ║")
	color.Blue.Println("╚════════════════════════════════════════════════════════════╝")
}

// ShowInstructions displays the usage instructions
func (u *UI) ShowInstructions() {
	color.Yellow.Println("\n📋 Select cleanup options:")
	color.Yellow.Println("   • Enter numbers separated by commas (e.g., 1,3,5)")
	color.Yellow.Println("   • Type 'all' to select all options")
	color.Yellow.Println("   • Type 'exit' or '8' to quit")
	color.Yellow.Println("\n🔧 Available Options:")
}

// ShowMenu displays the main menu with all options
func (u *UI) ShowMenu() {
	u.ClearScreen()
	u.ShowHeader("0.1.0") // TODO: Make version configurable
	u.ShowInstructions()

	options := []struct {
		icon  string
		text  string
		color func(a ...interface{}) string
	}{
		{"📦", "Clean npm cache", color.Green.Render},
		{"🧶", "Clean yarn cache", color.Green.Render},
		{"🐳", "Clean Docker cache", color.Green.Render},
		{"🪟", "Clean WinSxS temp files", color.Green.Render},
		{"🗑️", "Clean Windows temporary files", color.Green.Render},
		{"📝", "Clean Windows error reporting chunks", color.Green.Render},
		{"📊", "Show cache sizes", color.Cyan.Render},
		{"🚪", "Exit", color.Red.Render},
	}

	for i, opt := range options {
		fmt.Printf("  %s %s %s\n",
			color.Yellow.Sprintf("%d.", i+1),
			opt.icon,
			opt.color(opt.text))
	}

	color.Blue.Println("\n╔════════════════════════════════════════════════════════════╗")
	color.Blue.Println("║                    Make your choice                         ║")
	color.Blue.Println("╚════════════════════════════════════════════════════════════╝")
}

// ShowSelectedOptions displays the selected cleanup options
func (u *UI) ShowSelectedOptions(options []string) {
	color.Cyan.Println("\n🎯 Selected options:")
	for _, opt := range options {
		switch opt {
		case "1", "npm":
			color.Green.Println("  • npm cache")
		case "2", "yarn":
			color.Green.Println("  • yarn cache")
		case "3", "docker":
			color.Green.Println("  • Docker cache")
		case "4", "winsxs":
			color.Green.Println("  • WinSxS temp files")
		case "5", "wintemp":
			color.Green.Println("  • Windows temporary files")
		case "6", "winchunks":
			color.Green.Println("  • Windows error reporting chunks")
		case "7", "report":
			color.Cyan.Println("  • Cache size report")
		}
	}
	fmt.Println()
}

// ShowAdminWarning displays the administrator privileges warning
func (u *UI) ShowAdminWarning() {
	color.Red.Print("\n🔒 Administrator Privileges Required 🔒\n")
	color.Red.Println("========================================")
	color.Yellow.Println("⚠️  This tool requires administrator privileges to clean system caches.")
	color.Yellow.Println("Please run this tool as administrator.")
	color.Yellow.Println("\nPress Enter to exit...")
	u.reader.ReadBytes('\n')
	os.Exit(1)
}

// ShowError displays an error message
func (u *UI) ShowError(err error) {
	color.Red.Printf("[error] %v\n", err)
}

// ShowSuccess displays a success message
func (u *UI) ShowSuccess(msg string) {
	color.Green.Printf("\n✨ %s\n", msg)
}

// ShowWarning displays a warning message
func (u *UI) ShowWarning(msg string) {
	color.Yellow.Printf("\n⚠️  %s\n", msg)
}

// ShowInfo displays an info message
func (u *UI) ShowInfo(msg string) {
	color.Blue.Printf("\nℹ️  %s\n", msg)
}

// ShowCleanupStart displays the cleanup start message
func (u *UI) ShowCleanupStart() {
	color.Yellow.Println("🔄 Starting cleanup process...")
	fmt.Println()
}

// ShowCleanupComplete displays the cleanup completion message
func (u *UI) ShowCleanupComplete(errCount int) {
	if errCount > 0 {
		color.Red.Printf("\n⚠️  Clearance completed with %d error(s). Some operations may have failed.\n", errCount)
	} else {
		color.Green.Println("\n✨ Clearance finished successfully!")
	}
}

// ReadInput reads user input
func (u *UI) ReadInput() string {
	color.Yellow.Print("\n👉 Enter your choice: ")
	input, _ := u.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// WaitForEnter waits for the user to press Enter
func (u *UI) WaitForEnter() {
	color.Cyan.Print("\nPress Enter to continue...")
	u.reader.ReadBytes('\n')
}

// ShowCacheSizeReport displays the cache size report
func (u *UI) ShowCacheSizeReport(sizes map[string]string) {
	color.Blue.Println("\n📊 Cache Size Report")
	color.Blue.Println("===================")

	for name, size := range sizes {
		switch size {
		case "Not found", "N/A", "Not installed", "Docker not running":
			color.Yellow.Printf("%s: %s\n", name, size)
		case "Error", "Error getting size":
			color.Red.Printf("%s: %s\n", name, size)
		default:
			color.Green.Printf("%s: %s\n", name, size)
		}
	}
}
