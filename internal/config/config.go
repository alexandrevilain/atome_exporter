package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Atome holds all atome related config
type Atome struct {
	Username string
	Password string
}

// Config is the struct holding all app's config
type Config struct {
	Atome Atome
}

// LoadFromEnv returns the config populated from environement
func LoadFromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("ATOME_EXPORTER", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
