package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// Config houses stuff
type Config struct {
	ID   string
	Role string

	DiscoveryQuestion         string
	DiscoveryAnswer           string
	DiscoveryBroadcastAddress string
	DiscoveryBroadcastPort    string
	DiscoveryReceivePort      string

	WebsocketPort string

	VideoDevicesID int
	AudioDeviceID  int
}

var once = sync.Once{}
var configSingleton *Config = &Config{}

// GetConfig retrieves a singleton config object
func GetConfig() *Config {
	once.Do(func() {
		// on first load, read config file into the singleton variable
		configFile, err := os.Open("config.json")
		if err != nil {
			panic(err)
		}
		defer configFile.Close()

		configContents, err := ioutil.ReadAll(configFile)
		if err != nil {
			panic(err)
		}

		// unmarshal the config contents into the
		json.Unmarshal(configContents, configSingleton)
	})

	return configSingleton
}