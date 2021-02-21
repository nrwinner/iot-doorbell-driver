package webrtc

import (
	"doorbell-camera/src/entities"
	"doorbell-camera/src/modules/config"
	"doorbell-camera/src/modules/media"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

// TODO(nrwinner) logging
// TODO(nrwinner) better error handling all throughout this file

func handleNewCommand(controller *WebRTCController, command entities.Command) {
	// create a new uuid for the connection
	id := uuid.New().String()

	// sanity check: verify that the uuid doesn't already exist in map, regenerate until does not exist
	for {
		// if id doesn't exist in map, break out of loop
		if _, ok := controller.connections[id]; !ok {
			break
		}

		// replace id with new id and check again since it already exists
		id = uuid.New().String()
	}

	// retrieve singleton config information
	c := config.GetConfig()

	// respond to requester with id
	command.Client.SendCommand(entities.Command{
		Path: NEW_CONFIRM_COMMAND,
		Args: map[string]string{
			"id": id,
		},
		FromId:         c.ID,
		TargetDeviceId: command.FromId,
	})

	peerConnection := createPeerConnection()

	controller.connections[id] = peerConnection

	// set disconnect handler for PeerConnection
	controller.connections[id].OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		handleConnectionStateChange(state, id, controller.connections)
	})

	// set ice-candidate handler for PeerConnection
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		// transmit the new ice candidate via the signal servers
		sendIceCandidate(id, candidate, command.Client, command.FromId)
	})

	// attach media streams to PeerConnection
	attachMedia(peerConnection)

	// transmit offer via signaling server
	sendOffer(id, peerConnection, command.Client, command.FromId)
}

func handleAnswerCommand(controller *WebRTCController, command entities.Command) {
	// retrieve the peer id and offer from the command's Args array
	id := command.Args["id"]
	answerStr := command.Args["answer"]

	// retrieve the existing peer from the controller, err if does not exist
	peer := controller.connections[id]

	if peer == nil {
		panic("no peer with id " + id)
	}

	// unmarshal the answer string into a SessionDescription instance
	var answer webrtc.SessionDescription
	err := json.Unmarshal([]byte(answerStr), &answer)
	if err != nil {
		panic(err)
	}

	//set remote description to answer from client
	err = peer.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}
}

// Accept remote ICE candidates from signaling server and add them
func handleCandidateCommand(controller *WebRTCController, command entities.Command) {
	// fetch peer id and candidate from command's Args
	id := command.Args["id"]
	candidateStr := command.Args["candidate"]
	peer := controller.connections[id]

	if peer == nil {
		panic("no peer with id " + id)
	}

	// unmarshall the candidate string into an instance of ICECandidateInit
	var candidate webrtc.ICECandidateInit
	err := json.Unmarshal([]byte(candidateStr), &candidate)
	if err != nil {
		panic(err)
	}

	// add the decoded ICE candidate to the peer
	err = peer.AddICECandidate(candidate)
	if err != nil {
		panic(err)
	}
}

// After a disconnect event occurs, clean up local state
func handleConnectionStateChange(state webrtc.PeerConnectionState, id string, connections map[string]*webrtc.PeerConnection) {
	if state.String() == "disconnected" {
		// close the connection on disconnect
		// FIXME:NickW does WebRTC have an automatic reconnect capability that we're disabling here?
		peer := connections[id]
		for _, t := range peer.GetTransceivers() {
			println("Stopping", t.Kind().String())
			err := t.Stop()
			if err != nil {
				panic(err)
			}
		}
		err := peer.Close()
		if err != nil {
			// noop, couldn't close connection
			fmt.Println("Could not close connection on disconnect")
		}

		// remove the connection from stateful map
		delete(connections, id)
	}

	println("WebRTC State:", state.String())
}

// Create a new PeerConnection with media engine
func createPeerConnection() *webrtc.PeerConnection {
	mediaEngine := webrtc.MediaEngine{}
	codecSelector := media.GetCodecSelector()
	codecSelector.Populate(&mediaEngine)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	// create a new PeerConnection and store it in map, keyed by uuid
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	return peerConnection
}

// Sends an ICECandidate instance to client via signaling server
func sendIceCandidate(id string, candidate *webrtc.ICECandidate, client entities.Client, targetDeviceID string) {
	if candidate != nil {
		c := config.GetConfig()
		cBytes, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			panic(err)
		}

		client.SendCommand(entities.Command{
			Path: "webrtc/candidate",
			Args: map[string]string{
				"id":        id,
				"candidate": string(cBytes),
			},
			FromId:         c.ID,
			TargetDeviceId: targetDeviceID,
		})
	}
}

// Creates and sends an offer to client via signaling server
func sendOffer(id string, peer *webrtc.PeerConnection, client entities.Client, targetDeviceID string) {
	offer, err := peer.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	peer.SetLocalDescription(offer)

	offerStr, err := json.Marshal(offer)
	if err != nil {
		panic(err)
	}

	c := config.GetConfig()

	// send new offer to target
	client.SendCommand(entities.Command{
		Path: OFFER_COMMAND,
		Args: map[string]string{
			"id":    id,
			"offer": string(offerStr),
		},
		FromId:         c.ID,
		TargetDeviceId: targetDeviceID,
	})
}

func attachMedia(peer *webrtc.PeerConnection) {
	tracks := media.GetTracks()

	for _, track := range tracks {
		println("Adding", track.Kind().String())
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n",
				track.ID(), err)
		})

		_, err := peer.AddTransceiverFromTrack(track,
			webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)

		if err != nil {
			println("SHOULDNT BE HERE")
			panic(err)
		}
	}
}
