package calendarKaban

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/model"
	"github.com/aZaZloKeR/CalendarKaban/cmd/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

const (
	sessionName        = "kaban"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestId
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuth                  = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
		logger:       logrus.New(),
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestId)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handlerUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/sessions", s.handlerSessionCreate()).Methods(http.MethodPost)

	s.router.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"), //The url pointing to API definition
	))

	private := s.router.PathPrefix("/calendar").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/event", s.handleCalendarPage()).Methods(http.MethodPost)
}
func (s *server) handleCalendarPage() http.HandlerFunc {
	type request struct {
		Date        time.Time `json:"date"`
		TimeStart   time.Time `json:"timeStart"`
		TimeEnd     time.Time `json:"timeEnd"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		allRows := s.store.GetEventRepo().CountRows()
		event := model.Event{
			Date:      req.Date,
			TimeStart: req.TimeStart,
			TimeEnd:   req.TimeEnd,
			UserId:    r.Context().Value(ctxKeyUser).(*model.User).ID,
		}

		suitableRows, err := s.store.GetEventRepo().CountSuitableRows(event)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if suitableRows != allRows {
			s.error(w, r, http.StatusBadRequest, errors.New("выбранный промежуток времени уже занят на данное число"))
			return
		}

		if err = s.store.GetEventRepo().Create(&event); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

// @Accept json
// @Router /session [post]
func (s *server) handlerSessionCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.GetUserRepo().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Values["user_id"] = u.ID

		if err = s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}

}

// @Accept json
// @Router /users [post]
func (s *server) handlerUsersCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := model.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.GetUserRepo().Create(&u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := logrus.WithFields(logrus.Fields{
			"remoute_addr": r.RemoteAddr,
			"request_id":   r.Context().Value(ctxKeyRequestId),
		})
		logger.Infof("Begin handle request: %s, %s", r.Method, r.RequestURI)

		timeStart := time.Now()
		next.ServeHTTP(w, r)

		customWriter := &responseWriter{
			ResponseWriter: w,
			code:           http.StatusOK,
		}
		logger.Infof("Request executed { %v } with status code %v %s",
			time.Now().Sub(timeStart),
			customWriter.code,
			http.StatusText(customWriter.code),
		)
	})
}

func (s *server) setRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUUID := uuid.New().String()
		w.Header().Set("X-Request-ID", newUUID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestId, newUUID)))
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuth)
			return
		}
		user, err := s.store.GetUserRepo().FindById(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuth)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	})

}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
