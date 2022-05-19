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
	ps := mempubsub.NewPubSub(mempubsub.Config{BufSize: 10})
	defer ps.Close()
	errAccumulator := &testsubscribers.Accumulator{}
	ps.Register(errAccumulator)
	ps.Register(testsubscribers.Printer{})

	faceAcculator := &testsubscribers.Accumulator{}
	facePs := mempubsub.NewPubSub(mempubsub.Config{BufSize: 10})
	defer facePs.Close()
	facePs.Register(faceAcculator)
	facePs.Register(testsubscribers.Printer{})

	PhotoInterval = time.Millisecond
	MaxAttempts = 5
	DetectionCertaintyThreshold = 5.0

	photogrpher, err := NewFacePhotographer(&fakecapturer.Capturer{}, ps)
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
