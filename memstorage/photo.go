package memstorage

import (
	"sync"

	"github.com/areknoster/attendgo/domain"
)

var _ domain.PhotoStorage = (*PhotoStorage)(nil)

type PhotoStorage struct {
	m sync.Map
}

func NewPhotoStorage() *PhotoStorage {
	return &PhotoStorage{}
}

func (s *PhotoStorage) Create(Photo domain.Photo) error {
	_, ok := s.m.Load(Photo.Ref)
	if ok {
		return domain.ErrStorageAlreadyExits
	}
	s.m.Store(Photo.Ref, Photo)
	return nil
}

func (s *PhotoStorage) Get(id domain.PhotoRef) (domain.Photo, error) {
	Photo, ok := s.m.Load(id)
	if !ok {
		return domain.Photo{}, domain.ErrDoesntExist
	}
	domainPhoto, ok := Photo.(domain.Photo)
	if !ok {
		panic("wrong type of data stored under Photo ID")
	}
	return domainPhoto, nil
}

func (s *PhotoStorage) List() []domain.PhotoRef {
	ids := make([]domain.PhotoRef, 0)
	s.m.Range(func(key, value any) bool {
		ids = append(ids, key.(domain.PhotoRef))
		return true
	})
	return ids
}
