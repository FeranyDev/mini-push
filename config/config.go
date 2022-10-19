package config

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
	"os"
)

var Deploy *Config

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	TgBot string `yaml:"tg_bot"`
}

func init() {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Panicf("Error reading file: %v", err)
	}

	if err = yaml.Unmarshal(file, &Deploy); err != nil {
		log.Panicf("Error while unmarshaling config: %s", err)
	}

	log.Debug("Read config.yaml file: %s", Deploy)
}
