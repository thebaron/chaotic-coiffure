package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"server"`
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not load config file: %v", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
