package calendarKaban

import (
	"database/sql"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/config"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq" // ...
	"net/http"
)

func Start(config *config.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}
func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
