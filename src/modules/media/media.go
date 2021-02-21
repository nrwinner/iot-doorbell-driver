package media

import (
	"sync"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/prop"

	// side-effect needed to register the camera
	_ "github.com/pion/mediadevices/pkg/driver/camera"

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

// GetTracks retrieves a slice of video tracks for use with the WebRTC module
func GetTracks() []mediadevices.Track {
	tracksOnce.Do(func() {
		s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
			Video: func(c *mediadevices.MediaTrackConstraints) {
				c.Width = prop.Int(720)
				c.Height = prop.Int(480)
			},
			Codec: GetCodecSelector(),
		})

		if err != nil {
			println("are we here?")
			panic(err)
		}

		mediaSingleton.tracks = s.GetVideoTracks()
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

		mediaSingleton.codecSelector = mediadevices.NewCodecSelector(
			mediadevices.WithVideoEncoders(&videoEncoderParams),
		)
	})

	return mediaSingleton.codecSelector
}
