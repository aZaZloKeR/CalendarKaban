package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int    `json:"_"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (u *User) Sanitize() {
	u.Password = ""
}
