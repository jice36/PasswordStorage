package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Dbname     string `yaml:"dbname"`
	} `yaml:"database"`
}

func NewConfig(configPath string) (*Config, error) {
	config := new(Config)
	fmt.Println(configPath)
	fileConfig, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer fileConfig.Close()

	d := yaml.NewDecoder(fileConfig)
	err = d.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, err
}
