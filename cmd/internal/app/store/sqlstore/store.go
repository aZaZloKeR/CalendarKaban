package sqlstore

import (
	"database/sql"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/store"
)

type Store struct {
	db        *sql.DB
	userRepo  *UserRepository
	eventRepo *EventRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserRepo() store.UserRepository {
	if s.userRepo != nil {
		return s.userRepo
	}
	s.userRepo = &UserRepository{
		store: s,
	}
	return s.userRepo
}

func (s *Store) GetEventRepo() store.EventRepository {
	if s.eventRepo != nil {
		return s.eventRepo
	}
	s.eventRepo = &EventRepository{
		store: s,
	}
	return s.eventRepo
}
