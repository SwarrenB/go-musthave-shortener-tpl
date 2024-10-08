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
	serverAddress ServerAddress
	url           string
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

func (appConfig *Config) GetConfigURL() string {
	return appConfig.url
}

func CreateConfig() *Config {
	return &Config{
		serverAddress: ServerAddress{
			Host: `localhost`,
			Port: `8080`,
		},
		url: `http://localhost:8080/`,
	}

}

func CreateConfigWithFlags() *Config {
	devConfig := CreateConfig()
	flag.Var(&devConfig.serverAddress, "a", "server address {host:port}")
	flag.StringVar(&devConfig.url, "b", "", "URL address http://localhost:8080/{id}")
	flag.Parse()

	if devConfig.serverAddress.Host == "" || devConfig.serverAddress.Port == "" {
		devConfig.serverAddress.Set("http://localhost:8080/")
	}
	if devConfig.url == "" {
		devConfig.url = fmt.Sprintf("http://%s/", devConfig.serverAddress.String())
	}
	if !strings.HasSuffix(devConfig.url, "/") {
		devConfig.url += "/"
	}
	return devConfig
}
