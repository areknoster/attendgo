package fakecapturer

import (
	"bytes"
	"embed"
	"image/jpeg"
	"path"
	"time"

	"github.com/areknoster/attendgo/domain"
	"github.com/areknoster/attendgo/photo"
)

//go:embed photos
var photosDir embed.FS

const photosDirPath = "photos"

var photos []domain.Photo

func init() {
	entires, err := photosDir.ReadDir(photosDirPath)
	if err != nil {
		panic(err)
	}

	for _, e := range entires {
		if !e.Type().IsRegular() {
			continue
		}
		file, err := photosDir.ReadFile(path.Join(photosDirPath, e.Name()))
		if err != nil {
			panic(err)
		}
		img, err := jpeg.Decode(bytes.NewReader(file))
		if err != nil {
			panic(err)
		}
		photos = append(photos, domain.Photo{
			Img:  img,
			Date: time.Now(),
			Ref:  domain.NewPhotoRef(),
		})
	}
}

type Capturer struct {
	index int
}

var _ photo.Capturer = (*Capturer)(nil)

func (c *Capturer) Capture() domain.Photo {
	p := photos[c.index%len(photos)]
	c.index++
	return p
}
