package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) EncryptPassword() (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *User) Sanitize() {
	u.Password = ""
}
