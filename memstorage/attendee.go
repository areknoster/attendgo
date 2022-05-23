package memstorage

import (
	"sync"

	"github.com/areknoster/attendgo/domain"
)

var _ domain.AtendeeStorage = (*AtendeeStorage)(nil)

type AtendeeStorage struct {
	m sync.Map
}

func NewAtendeeStorage() *AtendeeStorage {
	return &AtendeeStorage{}
}

func (s *AtendeeStorage) Create(attendee domain.Attendee) error {
	_, ok := s.m.Load(attendee.ID)
	if ok {
		return domain.ErrStorageAlreadyExits
	}
	s.m.Store(attendee.ID, attendee)
	return nil
}

func (s *AtendeeStorage) Get(id domain.Attendee) (domain.Attendee, error) {
	attendee, ok := s.m.Load(id)
	if !ok {
		return domain.Attendee{}, domain.ErrDoesntExist
	}
	domainAttendee, ok := attendee.(domain.Attendee)
	if !ok {
		panic("wrong type of data stored under attendee ID")
	}
	return domainAttendee, nil
}

func (s *AtendeeStorage) List() []domain.Attendee {
	attendees := make([]domain.Attendee, 0)
	s.m.Range(func(key, value any) bool {
		attendees = append(attendees, value.(domain.Attendee))
		return true
	})
	return attendees
}
