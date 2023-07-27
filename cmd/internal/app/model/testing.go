package model

import "testing"

func TestUser(t *testing.T) *User {
	encPass, err := encryptString("azazloker")
	if err != nil {
		t.Fatal("cant encrypt password.", err)
	}

	return &User{
		Username: "azazloker",
		Email:    "azazloker@gmail.com",
		Password: encPass,
	}
}
