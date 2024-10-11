package config

import (
	"flag"
	"fmt"
	"strings"
)

type ServerAddress struct {
	Host string
	Port string
}

type Config struct {
	ServerAddress ServerAddress
	ShortURL      string
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
	port, _ := strings.CutSuffix(hp[1], "/")
	s.Port = port
	return nil
}

func CreateConfig() *Config {
	return &Config{
		ServerAddress: ServerAddress{
			Host: `localhost`,
			Port: `8080`,
		},
		ShortURL:   `http://localhost:8080/`,
		Vocabulary: make(map[string]string),
	}

}

func CreateConfigWithFlags() *Config {
	devConfig := CreateConfig()
	flag.Var(&devConfig.ServerAddress, "a", "server address {host:port}")
	flag.StringVar(&devConfig.ShortURL, "b", "", "URL address http://localhost:8080/{id}")
	flag.Parse()

	if devConfig.ServerAddress.Host == "" || devConfig.ServerAddress.Port == "" {
		devConfig.ServerAddress.Set("http://localhost:8080/")
	}
	if devConfig.ShortURL == "" {
		devConfig.ShortURL = fmt.Sprintf("http://%s/", devConfig.ServerAddress.String())
	}
	devConfig.ShortURL, _ = strings.CutSuffix(devConfig.ShortURL, "/")

	return devConfig
}
