package config

// Config holds the application configuration
type Config struct {
	Version string
	Commit  string
	Date    string
}

// NewConfig creates a new Config instance
func NewConfig(version, commit, date string) *Config {
	return &Config{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

// GetVersion returns the version information
func (c *Config) GetVersion() string {
	return c.Version
}

// GetCommit returns the commit hash
func (c *Config) GetCommit() string {
	return c.Commit
}

// GetDate returns the build date
func (c *Config) GetDate() string {
	return c.Date
}
