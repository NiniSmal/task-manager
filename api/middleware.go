package api

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
	"log"
	"net/http"
)

type Middleware struct {
	repo SessionRepository
}

func NewMiddleware(r SessionRepository) *Middleware {
	return &Middleware{
		repo: r,
	}
}

type SessionRepository interface {
	GetSession(ctx context.Context, sessionID uuid.UUID) (entity.User, error)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)

	})
}

func ResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("session_id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				HandlerError(w, entity.ErrNotAuthenticated)
				return
			}
			HandlerError(w, err)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value) //записанное значение Cookie при создании клиента записываем в sessionID
		if err != nil {
			HandlerError(w, err)
			return
		}

		user, err := m.repo.GetSession(ctx, sessionID) //в таблице связи  возвращаем нужный userID
		if err != nil {
			HandlerError(w, err)
			return
		}

		ctx = context.WithValue(ctx, "user", user) //вносим в context
		r = r.WithContext(ctx)                     //перезаписываем запрос с новым контестом, в который сохранили userID
		next.ServeHTTP(w, r)
	})
}
