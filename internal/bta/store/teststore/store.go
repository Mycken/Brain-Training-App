package teststore

import (
	"BrainTApp/internal/bta/store"
	"BrainTApp/internal/model/entity"
)

type Store struct {
	userRepository *UserRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*entity.User),
	}

	return s.userRepository
}
