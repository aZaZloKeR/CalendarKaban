package sqlstore

import (
	"database/sql"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/store"
)

type Store struct {
	db       *sql.DB
	userRepo *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) getUserRepo() *store.UserRepository {
	if s.userRepo != nil {
		return s.userRepo
	}
	s.userRepo = &UserRepository{
		store: s,
	}
	return s.userRepo
}
