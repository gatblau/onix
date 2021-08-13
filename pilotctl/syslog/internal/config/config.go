package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"5514"`
	}
	MongoDB struct {
		Host       string `yaml:"host" env-required:"true"`
		Port       string `yaml:"port"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		AuthDB     string `yaml:"auth_db" env-required:"true"`
		Database   string `yaml:"database" env-required:"true"`
		Collection string `yaml:"collection" env-required:"true"`
	} `yaml:"mongodb" env-required:"true"`
}

var instance *Config

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Printf("Info: %s", help)
			log.Fatal(err)
		}
	})
	return instance
}
