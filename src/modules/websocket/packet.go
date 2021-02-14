package websocket

import "doorbell-camera/src/entities"

const (
	COMMAND_PACKET = "command"
	ERROR_PACKET   = "error"
)

type Packet struct {
	PacketType string
}

type CommandPacket struct {
	Packet
	Command string
	Args    map[string]string
}

type ErrorPacket struct {
	Packet
	Error string
}

func CommandFromPacket(packet CommandPacket) entities.Command {
	return entities.Command{
		Path: packet.Command,
		Args: packet.Args,
	}
}

func PacketFromCommand(command entities.Command) CommandPacket {
	return CommandPacket {
		Packet: Packet{
			PacketType: COMMAND_PACKET,
		},
		Command: command.Path,
		Args: command.Args,
	}
}