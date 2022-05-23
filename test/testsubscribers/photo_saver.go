package testsubscribers

import (
	"image/jpeg"
	"log"
	"os"
	"path"

	"github.com/areknoster/attendgo/domain"
)

type PhotoSaver struct {
	Dir string
}

var _ domain.Subscriber = PhotoSaver{}

// Handle implements domain.Subscriber
func (ps PhotoSaver) Handle(ev domain.Event) {
	switch ev := ev.(type) {
	case domain.EventFacePhotoTaken:
		file, err := os.OpenFile(path.Join(ps.Dir, ev.Photo.Ref.String()), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
		if err != nil {
			log.Print("could not open image for write: ", err)
		}
		defer file.Close()
		err = jpeg.Encode(file, ev.Photo.Img, &jpeg.Options{
			Quality: 80,
		})
		if err != nil {
			log.Print("encode jped image: ", err)
		}
	}
}
