package config

import (
	"encoding/json"
	"github.com/kardianos/osext"
	"io/ioutil"
	"log"
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
		NewRelic struct {
			License string
			App     string
		}
		Rollbar struct {
			Token       string
			Environment string
		}
	}
}

var Config ConfigStruct

func LoadConfig(filePath string) {
	XP, _ := osext.ExecutableFolder()

	// Get the config file
	config_file, err := ioutil.ReadFile(XP + filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(config_file, &Config)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
}
