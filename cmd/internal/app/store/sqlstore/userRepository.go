package sqlstore

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	password, err := u.EncryptPassword()
	if err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO calendar.user (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		u.Username,
		u.Email,
		password,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, username, email, password FROM calendar.user where email = $1", email,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
	); err != nil {
		return nil, err
	}
	return u, nil
}
