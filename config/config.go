package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config represents the configuration information.
type ConfigStruct struct {
	Debug    bool
	Target   string
	Services struct {
		Mandrill struct {
			Key     string
			Account string
			Sender  struct {
				Name  string
				Email string
			}
		}
		Rollbar struct {
			Token       string
			Environment string
		}
	}
}

var Config ConfigStruct

func LoadConfig(filePath string) {
	// Get the config file
	config_file, err := ioutil.ReadFilefilePath)
	if err != nil {
		panic(err)
	}
	if  json.Unmarshal(config_file, &Config) != nil {
		panic(err)
	}
}

