package entities

type Command struct {
	Path           string
	Args           map[string]string
	Client         Client
	FromId         string
	TargetDeviceId string
}
