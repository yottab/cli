package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	DEFAULTE_CONTOROLLER = "controller.yottab.io:443"
	KEY_USER             = "yb-user"
	KEY_PASSWORD         = "yb-pass"
	KEY_TOKEN            = "token"
	KEY_HOST             = "host"
	APP_NAME             = "yb"
	CONFIG_NAME          = "config.json"
)

var (
	//ConfigManualAddress Get config file Path from the flag.
	ConfigManualAddress = ""

	// Find home directory for definition of archive folder
	// configPath == $HOME/.yb/
	configPath = filepath.Join(getHome(), "."+APP_NAME)
)

// Find home directory.
func getHome() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func isJSON(in []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(in, &js) == nil

}

//ensureConfigFile Ensure that a valid config file exists
func ensureConfigFile(filename string) {
	var err error
	var f *os.File
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		configRaw, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Panicf("Failed to read the config file: %v", err)
		}
		if isJSON(configRaw) {
			return
		}
		flags := os.O_WRONLY
		f, err = os.OpenFile(filename, flags, os.FileMode(0644))
	} else {
		err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		if err != nil {
			log.Panicf("Failed to create the config directory: %v", err)
		}
		flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		f, err = os.OpenFile(filename, flags, os.FileMode(0644))
		if err != nil {
			log.Panicf("Failed to create the config file: %v", err)
		}

	}
	defer f.Close()
	_, err = f.Write([]byte("{}"))
	if err != nil {
		log.Panicf("Failed to write to the config file: %v", err)
	}
}

//UpdateVarByConfigFile read config file into viper
func UpdateVarByConfigFile() {
	// read config either from Flag --config or Path "$HOME/.yb/config.json"
	var filename string
	if ConfigManualAddress != "" {
		// Use config file from the flag.
		filename = ConfigManualAddress
	} else {
		// Search config in home directory
		filename = filepath.Join(configPath, CONFIG_NAME)
	}
	ensureConfigFile(filename)
	viper.SetConfigFile(filename)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Failed to read the config file: %v", err)
	}
}

//ResetConfigFile wipes config file
func ResetConfigFile() (err error) {
	return viper.WriteConfig()
}
