package domain

import (
	"context"
	"image"
)

type Subscriber[E Event] interface {
	// Handle is called in serial manner
	Handle(ev E)
}

type SubscriberRegistry[E Event] interface {
	Register(s Subscriber[E])
}

type Publisher[E Event] interface {
	Publish(ev E)
}

type PubSub[E Event] interface {
	Publisher[E]
	SubscriberRegistry[E]
}

type Runner interface {
	Run(ctx context.Context) error
}

type EventKeyClicked struct {
	Glyph rune
}

func (e EventKeyClicked) String() string {
	return string(e.Glyph)
}

type EventIDInput string

func (e EventIDInput) String() string {
	return string(e)
}

type Photo struct {
	Img  image.Image
	Name string
}

type EventFacePhotoTaken struct {
	Photo Photo
}

func (e EventFacePhotoTaken) String() string {
	return e.Photo.Name
}

type Event interface {
	EventIDInput | EventKeyClicked | EventError | EventFacePhotoTaken
}
