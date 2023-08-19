package config

import (
	"os"
	"strings"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
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

	LocalMemory bool `yaml:"local_memory"`
}

func init() {

	config := "./config.yaml"
	environ := os.Environ()
	for _, e := range environ {
		if i := strings.Index(e, "="); i > 0 {
			key := e[:i]
			value := e[i+1:]
			switch key {
			case "MINI_PUSH_CONFIG":
				config = value
			}
		}
		file, err := os.ReadFile(config)
		if err != nil {
			log.Panicf("Error reading file: %v", err)
		}

		if err = yaml.Unmarshal(file, &Deploy); err != nil {
			log.Panicf("Error while unmarshaling config: %s", err)
		}

		log.Debug("Read config.yaml file: %s", Deploy)
	}
}
