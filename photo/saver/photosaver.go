package photosaver

import (
	"log"

	"github.com/areknoster/attendgo/domain"
)

var _ domain.Subscriber = (*PhotoSaver)(nil)

type PhotoSaver struct {
	storage domain.PhotoStorage
}

func New(s domain.PhotoStorage) *PhotoSaver{
	return &PhotoSaver{
		storage: s,
	}
}

// Handle implements domain.Subscriber
func (ps *PhotoSaver) Handle(ev domain.Event) {
	switch ev := ev.(type){
	case domain.EventFacePhotoTaken:
		err := ps.storage.Create(ev.Photo)
		if err != nil {
			log.Print("create photo in storage: ", err.Error())
			return
		}
	}
}
