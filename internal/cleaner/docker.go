package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DockerCleaner handles cleaning of Docker cache
type DockerCleaner struct {
	*BaseCleaner
}

// NewDockerCleaner creates a new DockerCleaner
func NewDockerCleaner() *DockerCleaner {
	return &DockerCleaner{
		BaseCleaner: NewBaseCleaner("docker"),
	}
}

// Clean performs the Docker cache cleaning operation
func (d *DockerCleaner) Clean(ctx context.Context) error {
	fmt.Println("[docker] Running Docker cleanup commands...")
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("[docker] Docker not found in PATH.")
		return fmt.Errorf("docker not found in PATH")
	}

	// Check if Docker daemon is running
	cmd := exec.CommandContext(ctx, "docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker daemon is not running")
	}

	commands := []struct {
		cmd  []string
		desc string
	}{
		{[]string{"docker", "system", "prune", "--all", "-f"}, "system prune"},
		{[]string{"docker", "volume", "prune", "-f"}, "volume prune"},
		{[]string{"docker", "builder", "prune", "--all", "-f"}, "builder prune"},
	}

	for _, c := range commands {
		cmd := exec.CommandContext(ctx, c.cmd[0], c.cmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("[docker] Failed to run %s: %v\n", c.desc, err)
			return fmt.Errorf("failed to run %s: %v", c.desc, err)
		}
	}

	fmt.Println("[docker] Docker cleanup completed.")
	return nil
}

// GetSize returns the size of Docker cache
func (d *DockerCleaner) GetSize(ctx context.Context) (string, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return "Not installed", nil
	}

	// Check if Docker daemon is running
	cmd := exec.CommandContext(ctx, "docker", "info")
	if err := cmd.Run(); err != nil {
		return "Docker not running", nil
	}

	cmd = exec.CommandContext(ctx, "docker", "system", "df", "--format", "{{.Size}}")
	output, err := cmd.Output()
	if err != nil {
		return "Error getting size", nil
	}
	return strings.TrimSpace(string(output)), nil
}
