package model_test

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/model"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestUser_EncryptString(t *testing.T) {
	u := model.TestUser(t)
	log.Printf("encpass for azazloker = %s", u.Password)
	assert.NotEmpty(t, u.Password)
}
