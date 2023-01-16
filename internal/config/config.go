package config

import (
	lg "beta/internal/domain/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Db struct {
		Xhost     string `yaml:"xhost"`
		Xport     string `yaml:"xport"`
		Xuser     string `yaml:"xuser"`
		Xpassword string `yaml:"xpassword"`
		Xdbname   string `yaml:"xdbname"`
	} `yaml:"db" env-prefix:"DB"`

	BaseConfig struct {
		HttpPort       string `yaml:"http_port"`
		JaegerEndpoint string `yaml:"jaeger_endpoint"`
		GammaEndpoint  string `yaml:"gamma_endpoint"`
	} `yaml:"base_config"`
}

var instance *Config
var once sync.Once

func GetConfig(logger lg.Factory) *Config {
	once.Do(func() {
		logger.Bg().Info("start reading configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			log.Fatal("can't read configuration")
		}
	})
	return instance
}
