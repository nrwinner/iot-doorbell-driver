package entities

type Controller interface {
	ParseCommand(client Client, command Command) int
	OnConnect(client Client)
	OnDisconnect(client Client)
}

type DefaultController struct {
	messageEventHandler    func(client Client, command Command)
	connectEventHandler    func(client Client)
	disconnectEventHandler func(client Client)
}

func (controller *DefaultController) New(
	messageEventHandler func(client Client, command Command),
	connectEventHandler func(client Client),
	disconnectEventHandler func(client Client),
) *DefaultController {
	controller.messageEventHandler = messageEventHandler
	controller.connectEventHandler = connectEventHandler
	controller.disconnectEventHandler = disconnectEventHandler

	return controller
}

func (controller *DefaultController) ParseCommand(client Client, command Command) int {
	if controller.messageEventHandler != nil {
		controller.messageEventHandler(client, command)
	}

	// pass through to next handler
	return 0
}

func (controller *DefaultController) OnConnect(client Client) {
	if controller.connectEventHandler != nil {
		controller.connectEventHandler(client)
	}
}

func (controller *DefaultController) OnDisconnect(client Client) {
	if controller.disconnectEventHandler != nil {
		controller.disconnectEventHandler(client)
	}
}
