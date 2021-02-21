package webrtc

import (
	"doorbell-camera/src/entities"
	"github.com/pion/webrtc/v3"
	"strings"
)

const (
	NEW_COMMAND         = "webrtc/new"
	NEW_CONFIRM_COMMAND = "webrtc/new-confirm"
	OFFER_COMMAND       = "webrtc/offer"
	ANSWER_COMMAND      = "webrtc/answer"
	CANDIDATE_COMMAND   = "webrtc/candidate"
)

type WebRTCConnection struct {
	connection *webrtc.PeerConnection
	socket entities.Client
}

type WebRTCController struct {
	connections map[string]*webrtc.PeerConnection
}

func (c *WebRTCController) ParseCommand(_ entities.Client, command entities.Command) int {
	if c.connections == nil {
		c.connections = make(map[string]*webrtc.PeerConnection)
	}

	if strings.HasPrefix(command.Path, "webrtc") {
		switch command.Path {
		case NEW_COMMAND:
			handleNewCommand(c, command)
			return 1
		case ANSWER_COMMAND:
			handleAnswerCommand(c, command)
			return 1
		case CANDIDATE_COMMAND:
			handleCandidateCommand(c, command)
			return 1
		}
	}

	return 0
}

func (c *WebRTCController) OnConnect(_ entities.Client) {}

func (c *WebRTCController) OnDisconnect(_ entities.Client) {}
