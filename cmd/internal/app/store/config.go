package store

type Config struct {
	databaseURL string `yaml:"databaseURL"`
}

func NewConfig() *Config {
	return &Config{}
}
