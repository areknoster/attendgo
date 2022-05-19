package domain

import (
	"context"
	"errors"
	"fmt"
	"image"
)

type Subscriber interface {
	// Handle is called in serial manner
	Handle(ev Event)
}

type SubscriberRegistry interface {
	Register(s Subscriber)
}

type Publisher interface {
	Publish(ev Event)
}

type PubSub interface {
	Publisher
	SubscriberRegistry
}

type Runner interface {
	Run(ctx context.Context) error
}

type EventKeyClicked struct {
	Glyph rune
}

func (e EventKeyClicked) String() string {
	return fmt.Sprintf("key clicked: %v", string(e.Glyph))
}

type ID string

type EventIDInput struct {
	ID ID
}

func (e EventIDInput) String() string {
	return fmt.Sprintf("id received: %s", string(e.ID))
}

type Photo struct {
	Img  image.Image
	Name string
}

type EventFacePhotoTaken struct {
	Photo Photo
}

func (e EventFacePhotoTaken) String() string {
	return fmt.Sprintf("face photo taken: %s", e.Photo.Name)
}

type Attendee struct {
	Photo Photo
	ID    ID
}

type AtendeeRegisteredEvent struct {
	Atendee Attendee
}

func (e AtendeeRegisteredEvent) String() string {
	return fmt.Sprintf("attendee registered id: %s photo: %s", e.Atendee.ID, e.Atendee.Photo.Name)
}

type AtendeeStorage interface {
	Create(attendee Attendee) error
	Get(id Attendee) (Attendee, error)
	List() []ID
}

var (
	ErrAttendeeDoesntExist         = errors.New("attendee does not exist")
	ErrStorageAttendeeAlreadyExits = errors.New("attendee already exists")
)

type Event interface{}
