package discovery

import (
	"net"
	"time"
)

const (
	broadcastPort    = "9999"
	receivePort      = "9998"
	broadcastAddress = "255.255.255.255"
	marco            = "marco"
	polo             = "polo"
)

// Attempt to locate the server on the network using a UDP packet.
// Broadcasts a question packet and listens for a corresponding answer packet
func Discover(r chan string) {
	pc, err := net.ListenPacket("udp4", ":"+receivePort)
	if err != nil {
		// assume no connection, noop
		return
	}

	// set read deadline, only try to read for one second
	pc.SetReadDeadline(time.Now().Add(time.Duration(time.Second)))

	defer pc.Close()

	// resolve UDP address
	udpAddress, err := net.ResolveUDPAddr("udp4", broadcastAddress+":"+broadcastPort)
	if err != nil {
		// UDP address could not be resolved
		panic(err)
	}

	// broadcast marco packet to address
	_, err = pc.WriteTo([]byte(marco), udpAddress)

	// Loop until one of the following conditions is met
	// 1. A response is located matching the answer
	// 2. A read timeout has occurred, indicating either a number of
	for {
		// listen for polo packet
		buffer := make([]byte, 1024)
		packetLength, serverAddress, err := pc.ReadFrom(buffer)
		if err != nil {
			// encountered read timeout, terminate
			r <- ""
			break
		} else {
			// we have a packet, check for match
			// verify that response is the answer to the question
			if string(buffer[:packetLength]) == polo {
				// respond with the address of the server
				r <- serverAddress.String()
				break
			}
		}
	}
}
