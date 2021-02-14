package websocket

import (
	"doorbell-camera/src/entities"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	Id     string
	Role   string
	socket *websocket.Conn
	mutex  sync.Mutex
}

func (s *Connection) GetId() string {
	return s.Id
}

func (s *Connection) GetRole() string {
	return s.Role
}

func (s *Connection) SendCommand(command entities.Command) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	packet := PacketFromCommand(command)
	err := s.socket.WriteJSON(packet)

	if err != nil {
		panic(err)
	}
}

func (s *Connection) SendError(error string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	packet := ErrorPacket{
		Packet: Packet{PacketType: ERROR_PACKET},
		Error:  error,
	}
	err := s.socket.WriteJSON(packet)

	if err != nil {
		panic(err)
	}
}
