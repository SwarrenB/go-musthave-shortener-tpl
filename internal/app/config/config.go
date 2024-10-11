package config

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type ServerAddress struct {
	Host string
	Port string
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	ShortURL      string `env:"BASE_URL"`
	Vocabulary    map[string]string
}

func (s ServerAddress) String() string {
	return s.Host + ":" + s.Port
}

func (s *ServerAddress) Set(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) != 2 {
		return fmt.Errorf("need address in a form host:port")
	}
	s.Host = hp[0]
	s.Port = hp[1]
	return nil
}

func CreateDefaultConfig() *Config {
	return &Config{
		ServerAddress: "localhost:8080",
		ShortURL:      `http://localhost:8080`,
		Vocabulary:    make(map[string]string),
	}

}

func CreateGeneralConfig() *Config {
	devConfig := CreateDefaultConfig()
	envConfig := Config{}
	flagsConfig := Config{}

	err := env.Parse(&envConfig)
	if err != nil {
		log.Fatal("Unable to parse environment config.")
	}

	flag.StringVar(&flagsConfig.ServerAddress, "a", "", "server address {host:port}")
	flag.StringVar(&flagsConfig.ShortURL, "b", "", "URL address http://localhost:8080/{id}")
	flag.Parse()

	if envConfig.ServerAddress != "" {
		devConfig.ServerAddress = envConfig.ServerAddress
	} else if flagsConfig.ServerAddress != "" {
		devConfig.ServerAddress = flagsConfig.ServerAddress
	}

	if envConfig.ShortURL != "" {
		devConfig.ShortURL = envConfig.ShortURL
	} else if flagsConfig.ShortURL != "" {
		devConfig.ShortURL = flagsConfig.ShortURL
	}

	return devConfig
}
