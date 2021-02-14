package entities

type Client interface {
	GetId() string
	GetRole() string
	SendCommand(command Command)
	SendError(error string)
}
