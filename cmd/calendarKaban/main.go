package main

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/calendarKaban"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/config"
	"log"
)

func main() {
	cfg := config.NewConfig()

	if err := calendarKaban.Start(cfg); err != nil {
		log.Fatal(err)
	}

}
