package domain

import "fmt"

type Attendee struct {
	PhotoRef PhotoRef
	ID       ID
}

type AtendeeRegisteredEvent struct {
	Atendee Attendee
}

func (e AtendeeRegisteredEvent) String() string {
	return fmt.Sprintf("attendee registered id: %s photo: %s", e.Atendee.ID, e.Atendee.PhotoRef)
}

type AtendeeStorage interface {
	Create(attendee Attendee) error
	Get(id Attendee) (Attendee, error)
	List() []Attendee
}
