package discovery

import (
	"doorbell-camera/src/modules/config"
	"fmt"
	"net"
	"time"
)

// Discover attempts to locate the server on the network using a UDP packet.
// Broadcasts a question packet and listens for a corresponding answer packet
func Discover(r chan string) {
	c := config.GetConfig()
	listenAddr, err := net.ResolveUDPAddr("udp", ":"+c.DiscoveryReceivePort)
	if err != nil {
		panic(err)
	}

	pc, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		// assume no connection, noop
		println("unable to resolve UDP address")
		return
	}

	// set read deadline, only try to read for n seconds
	pc.SetReadDeadline(time.Now().Add(time.Duration(5 * time.Second)))

	defer pc.Close()

	// resolve UDP address
	udpAddress, err := net.ResolveUDPAddr("udp4", c.DiscoveryBroadcastAddress+":"+c.DiscoveryBroadcastPort)
	if err != nil {
		// UDP address could not be resolved
		panic(err)
	}

	// broadcast marco packet to address
	_, err = pc.WriteTo([]byte(c.DiscoveryQuestion), udpAddress)

	// Loop until one of the following conditions is met
	// 1. A response is located matching the answer
	// 2. A read timeout has occurred, indicating either a number of
	for {
		// listen for polo packet
		buffer := make([]byte, 1024)
		println("reading...")
		packetLength, serverAddress, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			// encountered read timeout, terminate
			r <- ""
			break
		} else {
			// we have a packet, check for match
			// verify that response is the answer to the question
			if string(buffer[:packetLength]) == c.DiscoveryAnswer {
				// respond with the address of the server
				r <- serverAddress.String()
				break
			}
		}
	}
}
