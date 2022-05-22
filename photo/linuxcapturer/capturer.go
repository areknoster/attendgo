package linuxcapturer

import (
	"fmt"
	"log"
	"time"

	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/photo"
	"github.com/blackjack/webcam"
	"github.com/lamhai1401/mediadevices/pkg/frame"
)

const (
	FormatJPEG webcam.PixelFormat = 0x4A504547
	FormatYUYV webcam.PixelFormat = 0x56595559
)

type Decoder frame.Decoder

func NewYUYVDecoder() (Decoder, error) {
	return frame.NewDecoder(frame.FormatYUYV)
}

func Open(format webcam.PixelFormat, decoder Decoder) (*Capturer, error) {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		return nil, fmt.Errorf("open camera stream: %w", err)
	}

	size := cam.GetSupportedFrameSizes(format)[0]
	gotFormat, _, _, err := cam.SetImageFormat(format, size.MaxWidth, size.MaxHeight)
	if err != nil {
		return nil, fmt.Errorf("could not set image format: %w", err)
	}
	if gotFormat != format {
		return nil, fmt.Errorf("did not correctly set format, got %x instead", gotFormat)
	}
	err = cam.StartStreaming()
	if err != nil {
		log.Panic("start streaming: ", err)
	}

	return &Capturer{
		cam:     cam,
		decoder: decoder,
		witdth:  int(size.MaxWidth),
		height:  int(size.MaxHeight),
	}, nil
}

var _ photo.Capturer = (*Capturer)(nil)

type Capturer struct {
	cam            *webcam.Webcam
	decoder        Decoder
	witdth, height int
}

func (c *Capturer) Capture() domain.Photo {
	err := c.cam.WaitForFrame(uint32(time.Second))
	if err != nil {
		log.Panic("wait for frame: ", err)
	}
	frame, err := c.cam.ReadFrame()
	if err != nil {
		log.Panic("capture frame: ", err)
	}
	img, close, err := c.decoder.Decode(frame, c.witdth, c.height)

	if err != nil {
		log.Panic("decode frame: ", err)
	}
	defer close()

	return domain.Photo{
		Img:  img,
		Name: time.Now().Format(time.RFC3339) + ".jpg",
	}
}

func (c *Capturer) Close() error {
	return c.Close()
}
