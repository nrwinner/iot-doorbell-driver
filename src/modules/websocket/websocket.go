package websocket

import (
	"doorbell-camera/src/entities"
	"doorbell-camera/src/modules/config"

	"github.com/gorilla/websocket"
)

func Connect(serverAddress string, controllers []entities.Controller, lostConnection chan bool) {
	c := config.GetConfig()

	socket, _, err := websocket.DefaultDialer.Dial("ws://"+serverAddress+":"+c.WebsocketPort, nil)
	if err != nil {
		println(err.Error())
		// we were never able to make a connection, take same action
		// as though a connection were established and lost
		lostConnection <- true
		return
	}

	// create and transmit init packet
	initPacket := PacketFromCommand(entities.Command{
		Path: "system/init",
		Args: map[string]string{
			"id":   c.ID,
			"role": c.Role,
		},
	})

	err = socket.WriteJSON(initPacket)
	if err != nil {
		println(err.Error())
		// weren't able to write init packet, error to lostConnection to retry
		lostConnection <- true
		return
	}

	client := Connection{
		Id:     c.ID,
		Role:   c.Role,
		socket: socket,
	}

	// enter read loop and block
	for {
		var packet CommandPacket
		err := client.socket.ReadJSON(&packet)

		if err != nil {
			// read error, assume disconnect
			_ = client.socket.Close()
			lostConnection <- true
			return
		} else {
			// pass message to socket controller
			for _, controller := range controllers {
				command := CommandFromPacket(packet)
				command.Client = &client
				controller.ParseCommand(&client, command)
			}
		}

	}
}

func createResponder(client entities.Client) entities.Responder {
	return entities.Responder{
		Respond: func(command entities.Command) {
			// fetch client
			if client != nil {
				client.SendCommand(command)
			} else {
				panic("client doesn't exist")
			}
		},
		RespondError: func(error string) {
			if client != nil {
				client.SendError(error)
			} else {
				panic("client doesn't exist")
			}
		},
	}
}
