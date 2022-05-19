package memstorage

import (
	"github.com/areknoster/attendgo/domain"
	"sync"
)

var _ domain.AtendeeStorage = (*Storage)(nil)

type Storage struct {
	m sync.Map
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Create(attendee domain.Attendee) error {
	_, ok := s.m.Load(attendee.ID)
	if ok {
		return domain.ErrStorageAttendeeAlreadyExits
	}
	s.m.Store(attendee.ID, attendee)
	return nil
}

func (s *Storage) Get(id domain.Attendee) (domain.Attendee, error) {
	attendee, ok := s.m.Load(id)
	if !ok {
		return domain.Attendee{}, domain.ErrAttendeeDoesntExist
	}
	domainAttendee, ok := attendee.(domain.Attendee)
	if !ok {
		panic("wrong type of data stored under attendee ID")
	}
	return domainAttendee, nil
}

func (s *Storage) List() []domain.ID {
	ids := make([]domain.ID, 0)
	s.m.Range(func(key, value any) bool {
		ids = append(ids, key.(domain.ID))
		return true
	})
	return ids
}
