package entities

type Responder struct {
	Respond func(command Command)
	RespondError func(error string)
}
