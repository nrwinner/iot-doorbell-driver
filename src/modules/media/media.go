package media

import (
	"sync"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/prop"

	// side-effect needed to register the microphone
	_ "github.com/pion/mediadevices/pkg/driver/audiotest"

	// side-effect needed to register the camera
	_ "github.com/pion/mediadevices/pkg/driver/camera"

	// load the opus codec as our audio encoder
	audioEncoder "github.com/pion/mediadevices/pkg/codec/opus"

	// load the mmal codec as our video encoder (uses rpi hardware encoding)
	videoEncoder "github.com/pion/mediadevices/pkg/codec/mmal"
)

var tracksOnce = sync.Once{}
var codecOnce = sync.Once{}

type media struct {
	tracks        []mediadevices.Track
	codecSelector *mediadevices.CodecSelector
}

var mediaSingleton *media = &media{}

// GetTracks retrieves a slice of audio/video tracks for use with the WebRTC module
func GetTracks() []mediadevices.Track {
	tracksOnce.Do(func() {
		s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
			Video: func(c *mediadevices.MediaTrackConstraints) {
				c.Width = prop.Int(1280)
				c.Height = prop.Int(720)
			},
			Audio: func(c *mediadevices.MediaTrackConstraints) {},
			Codec: GetCodecSelector(),
		})

		if err != nil {
			println("are we here?")
			panic(err)
		}

		mediaSingleton.tracks = s.GetTracks()
	})

	return mediaSingleton.tracks
}

// GetCodecSelector creates/retrieves a selector object for use with the WebRTC module
func GetCodecSelector() *mediadevices.CodecSelector {
	codecOnce.Do(func() {
		videoEncoderParams, err := videoEncoder.NewParams()
		if err != nil {
			panic(err)
		}
		videoEncoderParams.BitRate = 500_000 // 500kbps

		// create params for new audio encoder
		audioEncoderParams, err := audioEncoder.NewParams()
		if err != nil {
			panic(err)
		}

		mediaSingleton.codecSelector = mediadevices.NewCodecSelector(
			mediadevices.WithVideoEncoders(&videoEncoderParams),
			mediadevices.WithAudioEncoders(&audioEncoderParams),
		)
	})

	return mediaSingleton.codecSelector
}
