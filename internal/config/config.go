package config

import (
	"flag"
	"fmt"
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
}

func Parse() *Config {
	c := &Config{}
	flag.IntVar(&c.Port, "port", 9080, "Listen port")
	flag.IntVar(&c.Port, "p", 9080, "Listen port (shorthand)")
	flag.StringVar(&c.Listen, "listen", "127.0.0.1", "Bind address (use 0.0.0.0 for remote access)")
	flag.StringVar(&c.Theme, "theme", "github", "Default CSS theme")
	flag.StringVar(&c.PlantUMLServer, "plantuml-server", "https://www.plantuml.com/plantuml", "PlantUML server URL")
	flag.BoolVar(&c.Verbose, "verbose", false, "Enable debug logging")
	flag.BoolVar(&c.Verbose, "v", false, "Enable debug logging (shorthand)")
	flag.StringVar(&c.AuthToken, "auth-token", "", "Authentication token (required in remote mode)")
	flag.IntVar(&c.AuthCookieMaxAge, "auth-cookie-max-age", 30, "Authentication cookie max age in days")
	flag.StringVar(&c.AccessLog, "access-log", "", "Access log file path")
	flag.IntVar(&c.AccessLogMaxSize, "access-log-max-size", 100, "Access log max size in MB before rotation")
	flag.IntVar(&c.AccessLogMaxBack, "access-log-max-backups", 3, "Max number of old access log files to retain")
	flag.IntVar(&c.AccessLogMaxAge, "access-log-max-age", 28, "Max number of days to retain old access log files")
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
