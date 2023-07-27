package apiserver

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/store"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Port     string        `yaml:"port"`
	LogLevel string        `yaml:"logLevel"`
	Store    *store.Config `yaml:"store"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
