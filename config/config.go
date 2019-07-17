package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config stores all the information in the configuration file
type Config struct {
	CheckForTODOs  bool          `yaml:"checkForTODOs"`
	CheckForPrints bool          `yaml:"checkForPrints"`
	Github         *GithubConfig `yaml:"github"`
}

// GithubConfig stores the configuration related to Github
type GithubConfig struct {
	Branch              string `yaml:"Branch"`
	Base                string `yaml:"defaultBase"`
	MaintainerCanModify bool   `yaml:"defaultMaintainerCanModify"`
	Organization        string `yaml:"organization"`
	Repository          string `yaml:"repository"`
	Username            string `yaml:"username"`
	IsFork              bool   `yaml:"isFork"`
}

// GetConfig returns the config struct given a path to the config file
func GetConfig(path string) (*Config, error) {
	config := &Config{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	yamlData, err := ioutil.ReadFile(homedir + "/InternLinter/config.yaml")
	err = yaml.UnmarshalStrict(yamlData, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
