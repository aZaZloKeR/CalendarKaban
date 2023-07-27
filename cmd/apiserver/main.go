package main

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/apiserver"
	"log"
)

func main() {
	cfg := apiserver.GetConfig()

	if err := apiserver.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
