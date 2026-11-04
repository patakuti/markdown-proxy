package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Port             int
	Listen           string
	Theme            string
	PlantUMLServer   string
	Verbose          bool
	AuthToken        string
	AuthCookieMaxAge int
	AccessLog        string
	AccessLogMaxSize int
	AccessLogMaxBack int
	AccessLogMaxAge  int
	Configure        bool
	NoBrowser        bool
}

// configFilePathFunc returns the path to the configuration file.
// It is a variable so tests can override it.
var configFilePathFunc = defaultConfigFilePath

func defaultConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "markdown-proxy", "config.json"), nil
}

// fileConfig holds the settings that can be persisted in the config file.
type fileConfig struct {
	PlantUMLServer string `json:"plantuml-server,omitempty"`
	Theme          string `json:"theme,omitempty"`
	Port           int    `json:"port,omitempty"`
	Listen         string `json:"listen,omitempty"`
	Verbose        bool   `json:"verbose,omitempty"`
	NoBrowser      bool   `json:"no-browser,omitempty"`
}

// loadConfigFile reads the config file and returns the parsed values.
// Returns a zero-value fileConfig if the file does not exist.
func loadConfigFile() fileConfig {
	var fc fileConfig
	path, err := configFilePathFunc()
	if err != nil {
		return fc
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fc
	}
	_ = json.Unmarshal(data, &fc)
	return fc
}

func Parse() *Config {
	c := &Config{}

	// Load config file values as defaults
	fc := loadConfigFile()
	defaultPort := 9080
	if fc.Port != 0 {
		defaultPort = fc.Port
	}
	defaultListen := "127.0.0.1"
	if fc.Listen != "" {
		defaultListen = fc.Listen
	}
	defaultTheme := "github"
	if fc.Theme != "" {
		defaultTheme = fc.Theme
	}
	defaultPlantUML := ""
	if fc.PlantUMLServer != "" {
		defaultPlantUML = fc.PlantUMLServer
	}
	defaultVerbose := fc.Verbose
	defaultNoBrowser := fc.NoBrowser

	flag.IntVar(&c.Port, "port", defaultPort, "Listen port")
	flag.IntVar(&c.Port, "p", defaultPort, "Listen port (shorthand)")
	flag.StringVar(&c.Listen, "listen", defaultListen, "Bind address (use 0.0.0.0 for remote access)")
	flag.StringVar(&c.Theme, "theme", defaultTheme, "Default CSS theme")
	flag.StringVar(&c.PlantUMLServer, "plantuml-server", defaultPlantUML, "PlantUML server URL (disabled by default; e.g. https://www.plantuml.com/plantuml)")
	flag.BoolVar(&c.Verbose, "verbose", defaultVerbose, "Enable debug logging")
	flag.BoolVar(&c.Verbose, "v", defaultVerbose, "Enable debug logging (shorthand)")
	flag.BoolVar(&c.NoBrowser, "no-browser", defaultNoBrowser, "Do not open a browser; print the URL and run the server in the foreground")
	flag.StringVar(&c.AuthToken, "auth-token", "", "Authentication token (required in remote mode)")
	flag.IntVar(&c.AuthCookieMaxAge, "auth-cookie-max-age", 30, "Authentication cookie max age in days")
	flag.StringVar(&c.AccessLog, "access-log", "", "Access log file path")
	flag.IntVar(&c.AccessLogMaxSize, "access-log-max-size", 100, "Access log max size in MB before rotation")
	flag.IntVar(&c.AccessLogMaxBack, "access-log-max-backups", 3, "Max number of old access log files to retain")
	flag.IntVar(&c.AccessLogMaxAge, "access-log-max-age", 28, "Max number of days to retain old access log files")
	flag.BoolVar(&c.Configure, "configure", false, "Interactively create configuration file")
	flag.Parse()
	return c
}

// IsRemoteMode returns true when the server is configured to accept
// connections from outside localhost.
func (c *Config) IsRemoteMode() bool {
	return c.Listen != "127.0.0.1" && c.Listen != "localhost"
}

// Validate checks the configuration for consistency.
// Returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if c.IsRemoteMode() && c.AuthToken == "" {
		return fmt.Errorf("--auth-token is required when --listen is not localhost (remote mode)")
	}
	return nil
}
