package sqlstore

import (
	"BrainTApp/internal/bta/store"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db               *sql.DB
	userRepository   *UserRepository
	resultRepository *ResultRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
func (s *Store) Result() store.ResultRepository {
	if s.resultRepository != nil {
		return s.resultRepository
	}

	s.resultRepository = &ResultRepository{
		store: s,
	}

	return s.resultRepository
}

//store.User.Create()
