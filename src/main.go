package main

import (
	"doorbell-camera/src/entities"
	"doorbell-camera/src/modules/config"
	"doorbell-camera/src/modules/discovery"
	"doorbell-camera/src/modules/websocket"
	"strings"
	"time"
)

func main() {
	hasConnection := false
	serverChan := make(chan string)
	connectionLost := make(chan bool)

	// request config to force read from file on start-up instead of JIT
	config.GetConfig()

	// main application loop - discover server -> do work -> lose connection -> repeat
	for {
		if hasConnection {
			// wait for connectionLost to signal and flip flag
			<-connectionLost
			hasConnection = false
			time.Sleep(5 * time.Second)
		} else {
			// we don't have a connection, attempt to discover first
			go discovery.Discover(serverChan)
			server := <-serverChan

			if server != "" {
				// connection is found
				hasConnection = true
				addressParts := strings.Split(server, ":")
				if len(addressParts) != 2 {
					panic("unexpected format for server address - " + server)
				}
				go start(addressParts[0], connectionLost)
			} else {
				// connection is not found
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// Takes the address of the server and makes appropriate connections to various servers
// Writes to lostConnection when a connection could not be established or is lost after being established
func start(server string, lostConnection chan bool) {
	go websocket.Connect(server, []entities.Controller{}, lostConnection)
}
