package photo_test

import (
	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/mempubsub"
	. "github.com/areknoster/attendgo/photo"
	"github.com/areknoster/attendgo/photo/fakecapturer"
	"github.com/areknoster/attendgo/test/testsubscribers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPhotographer(t *testing.T) {
	errPs := mempubsub.NewPubSub[domain.EventError](mempubsub.Config{BufSize: 10})
	defer errPs.Close()
	errAccumulator := &testsubscribers.Accumulator[domain.EventError]{}
	errPs.Register(errAccumulator)
	errPs.Register(testsubscribers.Printer[domain.EventError]{})

	faceAcculator := &testsubscribers.Accumulator[domain.EventFacePhotoTaken]{}
	facePs := mempubsub.NewPubSub[domain.EventFacePhotoTaken](mempubsub.Config{BufSize: 10})
	defer facePs.Close()
	facePs.Register(faceAcculator)
	facePs.Register(testsubscribers.Printer[domain.EventFacePhotoTaken]{})

	PhotoInterval = time.Millisecond
	MaxAttempts = 5
	DetectionCertaintyThreshold = 5.0

	photogrpher, err := NewFacePhotographer(&fakecapturer.Capturer{}, errPs, facePs)
	require.NoError(t, err)
	photogrpher.StartCapturing()
	time.Sleep(time.Second)

	assert.EqualValues(t, []domain.EventError{
		domain.ErrNoFace,
		domain.ErrTooManyFaces,
		domain.ErrTooManyFaces,
		domain.ErrNoFace,
	}, errAccumulator.Received)

	assert.Len(t, faceAcculator.Received, 1)
}
