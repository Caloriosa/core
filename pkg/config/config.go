package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	MongoConnection string `yaml:"MongoConnection"`
	MongoDatabase   string `yaml:"MongoDatabase"`
	Email           struct {
		SmtpHost     string `yaml:"SmtpHost"`
		SmtpPort     int    `yaml:"SmtpPort"`
		SmtpUser     string `yaml:"SmtpUser"`
		SmtpPassword string `yaml:"SmtpPassword"`
		EmailFrom    string `yaml:"EmailFrom"`
	} `yaml:"Email"`

	Dev struct {
		TestEmailTo string `yaml:"TestEmailTo"`
	} `yaml:"Dev"`

	AppTokens []struct {
		App   string `yaml:"App"`
		Token string `yaml:"Token"`
	} `yaml:"AppTokens"`
}

var LoadedConfig *Config

func LoadConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	LoadedConfig = &Config{}
	err = yaml.Unmarshal(yamlFile, &LoadedConfig)

	return err
}
