package config

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
	"os"
)

var Deploy *Config

type RedisDB struct {
	Addr string `yaml:"addr"`

	User string `yaml:"user"`
	Pass string `yaml:"pass"`

	Db int `yaml:"db"`
}

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Debug bool `yaml:"debug"`

	TgBot string `yaml:"tg_bot"`

	Redis RedisDB `yaml:"redis"`
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
