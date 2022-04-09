package practical_go

import "charliemcelfresh/practical_go/internal/config"

type Config struct {
	*config.Config
}

func NewConfig() *Config {
	return &Config{
		Config: config.NewConfig(),
	}
}
