package config

import "flag"

type Config struct {
	Port           int
	Theme          string
	PlantUMLServer string
}

func Parse() *Config {
	c := &Config{}
	flag.IntVar(&c.Port, "port", 9080, "Listen port")
	flag.IntVar(&c.Port, "p", 9080, "Listen port (shorthand)")
	flag.StringVar(&c.Theme, "theme", "github", "Default CSS theme")
	flag.StringVar(&c.PlantUMLServer, "plantuml-server", "https://www.plantuml.com/plantuml", "PlantUML server URL")
	flag.Parse()
	return c
}
