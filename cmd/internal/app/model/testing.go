package model

import "testing"

func TestUser(t *testing.T) *User {
	user := User{
		Username: "azazloker",
		Email:    "azazloker@gmail.com",
		Password: "azazloker",
	}
	encPass, err := user.EncryptPassword()
	if err != nil {
		t.Fatal("cant encrypt password.", err)
	}
	user.Password = encPass
	return &user
}
