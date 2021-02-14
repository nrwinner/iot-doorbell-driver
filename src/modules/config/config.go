package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type config struct {
	Id   string
	Role string

	DiscoveryQuestion         string
	DiscoveryAnswer           string
	DiscoveryBroadcastAddress string
	DiscoveryBroadcastPort    string
	DiscoveryReceivePort      string

	WebsocketPort string
}

var singleton *config

func GetConfig() *config {
	once := sync.Once{}
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

		json.Unmarshal(configContents, &singleton)
	})

	return singleton
}
