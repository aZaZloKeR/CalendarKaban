package sqlstore

import "github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/model"

type EventRepository struct {
	store *Store
}

func (r *EventRepository) Create(e *model.Event) error {
	return r.store.db.QueryRow(
		"INSERT INTO calendar.event (date, time_start, time_end, name, description, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		e.Date,
		e.TimeStart,
		e.TimeEnd,
		e.Name,
		e.Description,
		e.UserId,
	).Scan(&e.ID)
}

func (r *EventRepository) FindById(id int) (*model.Event, error) {
	e := &model.Event{}

	if err := r.store.db.QueryRow(
		"SELECT id, date, time_start, time_end, name, description, user_id FROM calendar.event WHERE id = $1", id,
	).Scan(
		&e.ID,
		&e.Date,
		&e.TimeStart,
		&e.TimeEnd,
		&e.Name,
		&e.Description,
		&e.UserId,
	); err != nil {
		return nil, err
	}
	return e, nil
}

func (r *EventRepository) CountSuitableRows(event model.Event) (int, error) {
	var count int
	if err := r.store.db.QueryRow(
		"SELECT count(id) FROM calendar.event WHERE user_id = $1 AND date = $2 AND ($3 < time_start AND $4 <= time_start) OR ($3 >= time_end AND $4 > time_end)",
		event.UserId, event.Date, event.TimeStart, event.TimeEnd,
	).Scan(count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *EventRepository) CountRows() int {
	var count int
	r.store.db.QueryRow("SELECT count(id) FROM calendar.event")
	return count
}
