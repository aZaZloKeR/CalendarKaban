package store

import "github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/model"

type UserRepository interface {
	Create(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}
