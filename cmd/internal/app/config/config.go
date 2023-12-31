package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	BindAddr    string `yaml:"bindAddr"`
	LogLevel    string `yaml:"logLevel"`
	DatabaseURL string `yaml:"databaseURL"`
	SessionKey  string `yaml:"sessionKey"`
}

var instance *Config
var once sync.Once

func NewConfig() *Config {
	once.Do(func() {
		log.Print("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("C:/Users/azazl/GolandProjects/CalendarKaban/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
