package main

import (
	"doorbell-camera/src/modules/config"
	"doorbell-camera/src/modules/discovery"
	"time"
)

func main() {
	hasConnection := false
	serverChan := make(chan string)
	connectionLost := make(chan bool)

	// request config to force read from file
	config.GetConfig()

	// main application loop - discover server -> do work -> lose connection -> repeat
	for {
		if hasConnection {
			// wait for connection lost to signal and flip flag
			<-connectionLost
			hasConnection = false
		} else {
			// we don't have a connection, attempt to discover first then make connection
			go discovery.Discover(serverChan)
			server := <-serverChan

			if server != "" {
				// connection is found
				hasConnection = true
				go start(server, connectionLost)
			} else {
				// connection is not found
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// Takes the address of the server and makes appropriate connections to various servers
// Writes to lostConnection when a connection could not be established or is lost after being established
func start(_ string, lostConnection chan bool) {
	// TODO:NickW implement websocket server here
}
