package domain

import (
	"fmt"
	"image"
	"time"

	"github.com/google/uuid"
)

func NewPhotoRef() PhotoRef {
	return PhotoRef(uuid.New())
}

type PhotoRef uuid.UUID

func (pr PhotoRef) String() string {
	return uuid.UUID(pr).String()
}

func ParsePhotoRef(r string) (PhotoRef, error) {
	id, err := uuid.Parse(r)
	if err != nil {
		return PhotoRef{}, err
	}
	return PhotoRef(id), nil
}

type Photo struct {
	Img  image.Image
	Date time.Time
	Ref  PhotoRef
}

type EventFacePhotoTaken struct {
	Photo Photo
}

func (e EventFacePhotoTaken) String() string {
	return fmt.Sprintf("face photo taken at %s: %s", e.Photo.Date.String(), e.Photo.Ref)
}

type PhotoStorage interface {
	Create(photo Photo) error
	Get(id PhotoRef) (Photo, error)
	List() []PhotoRef
}
