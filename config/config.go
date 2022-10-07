package config

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Version string `yaml:"version"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

func init() {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Panicf("Error reading file: %v", err)
	}

	var config Config

	if err = yaml.Unmarshal(file, &config); err != nil {
		log.Panicf("Error while unmarshaling config: %s", err)
	}

	log.Debug("Read config.yaml file: %s", config)
}
