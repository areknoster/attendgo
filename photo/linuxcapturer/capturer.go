package linuxcapturer

import (
	"bytes"
	"fmt"
	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/photo"
	"github.com/blackjack/webcam"
	"image/jpeg"
	"log"
	"time"
)

const FormatJPEG = 0x4A504547

func Open() (*Capturer, error) {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		return nil, fmt.Errorf("open camera stream: %w", err)
	}
	size := cam.GetSupportedFrameSizes(FormatJPEG)[0]
	gotFormat, _, _, err := cam.SetImageFormat(FormatJPEG, size.MaxWidth, size.MaxHeight)
	if err != nil {
		return nil, fmt.Errorf("could not set image format: %w", err)
	}
	if gotFormat != FormatJPEG {
		return nil, fmt.Errorf("did not correctly set JPEG format, got %x instead", gotFormat)
	}

	return &Capturer{
		cam: cam,
	}, nil
}

var _ photo.Capturer = (*Capturer)(nil)

type Capturer struct {
	cam *webcam.Webcam
}

func (c *Capturer) Capture() domain.Photo {
	frame, err := c.cam.ReadFrame()
	if err != nil {
		log.Fatal("capture frame: ", err)
	}
	img, err := jpeg.Decode(bytes.NewReader(frame))
	if err != nil {
		log.Fatal("decode JPEG frame: ", err)
	}
	return domain.Photo{
		Img:  img,
		Name: time.Now().Format(time.RFC3339) + ".jpg",
	}
}

func (c *Capturer) Close() error {
	return c.Close()
}
