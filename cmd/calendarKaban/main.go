package main

import (
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/calendarKaban"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/config"
	"log"
)

// @title Календарь
// @version 1.0
// description Какой то спецефичный календарь

// @host localhost:8081
// @BasePath /

func main() {
	cfg := config.NewConfig()

	if err := calendarKaban.Start(cfg); err != nil {
		log.Fatal(err)
	}

}
