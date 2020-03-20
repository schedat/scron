package server

import (
	"io"

	"github.com/spf13/viper"
)

// job contains information about a job: ID, name and how to
// execute it with Program add Arguments
type job struct {
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Program   string `mapstructure:"program"`
	Arguments string `mapstructure:"arguments"`
	Enabled   bool   `mapstructure:"enabled"`
}

// Config contains configuration for Scheduler
type Config struct {
	Host string `mapstructure:"host"`
	Jobs []job  `mapstructure:"jobs"`
}

// Parse parses config
func Parse(reader io.Reader) (*Config, error) {
	viper.SetConfigType("yaml")

	var cfg Config

	err := viper.ReadConfig(reader)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
