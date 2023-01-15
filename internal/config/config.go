package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	lg "gitlab.satel.eyevox.ru/satel_vks/jaeger_tracer/log"
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
		Time           string `yaml:"time_to_live"`
		HttpPort       string `yaml:"http_port"`
		JaegerEndpoint string `yaml:"jaeger_endpoint"`
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
