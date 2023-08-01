package model

import "time"

type Event struct {
	ID          int       `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	TimeStart   time.Time `json:"timeStart"`
	TimeEnd     time.Time `json:"timeEnd"`
	UserId      int       `json:"-"`
}
