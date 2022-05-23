package domain

import (
	"context"
	"errors"
	"fmt"
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

var (
	ErrDoesntExist         = errors.New("attendee does not exist")
	ErrStorageAlreadyExits = errors.New("attendee already exists")
)

type Event interface{}
