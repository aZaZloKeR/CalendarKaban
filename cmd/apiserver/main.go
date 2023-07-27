package main

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/apiserver"
	"log"
)

func main() {
	cfg := apiserver.GetConfig()
	s := apiserver.New(cfg)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
