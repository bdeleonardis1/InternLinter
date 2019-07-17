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
	Title               string
}

// GetConfig returns the config struct given a path to the config file
func GetConfig(path string) (*Config, error) {
	config := &Config{}
	if path == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = homedir + "/InternLinter/config.yaml"
	}
	yamlData, err := ioutil.ReadFile(path)
	err = yaml.UnmarshalStrict(yamlData, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetArgs returns a map of command line arguments.
func GetArgs() map[string]string {
	commandlineArgs := make(map[string]string)
	lastKey := ""
	for _, arg := range os.Args[1:] {
		if lastKey != "" {
			commandlineArgs[lastKey] = arg
			lastKey = ""
		} else {
			lastKey = arg
		}
	}

	return commandlineArgs
}
