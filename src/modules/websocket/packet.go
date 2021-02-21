package websocket

import "doorbell-camera/src/entities"

const (
	COMMAND_PACKET = "command"
	ERROR_PACKET   = "error"
)

type Packet struct {
	PacketType     string
	FromId         string
	TargetDeviceId string
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
		FromId: packet.FromId,
		TargetDeviceId: packet.TargetDeviceId,
	}
}

func PacketFromCommand(command entities.Command) CommandPacket {
	return CommandPacket{
		Packet: Packet{
			PacketType:     COMMAND_PACKET,
			FromId:         command.FromId,
			TargetDeviceId: command.TargetDeviceId,
		},
		Command: command.Path,
		Args:    command.Args,
	}
}
