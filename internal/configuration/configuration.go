package configuration

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	originalFile string        `yaml:"-"`
	BindAddress  string        `yaml:"bind_address"`
	Auth         string        `yaml:"auth"`
	Replications []Replication `yaml:"replications"`
}

type Replication struct {
	Address string
	Auth    string
}

func FromFile(file string) (*AppConfig, error) {
	cont, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var app AppConfig
	if err := yaml.Unmarshal(cont, &app); err != nil {
		return nil, err
	}

	app.originalFile = file
	return &app, nil
}

func (a *AppConfig) ToFile(name string) error {
	fileToWrite := a.originalFile
	if name != "" {
		fileToWrite = name
	}

	y, err := yaml.Marshal(*a)
	if err != nil {
		return err
	}

	return os.WriteFile(fileToWrite, y, 0644)
}

func LoadDefaults() *AppConfig {
	return &AppConfig{
		BindAddress: "localhost:7000",
	}
}
