package store

import "github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/model"

type UserRepository interface {
	Create(u *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindById(id int) (*model.User, error)
}

type EventRepository interface {
	Create(u *model.Event) error
	FindById(id int) (*model.Event, error)
	CountSuitableRows(event model.Event) (int, error)
	CountRows() int
}
