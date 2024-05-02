package api

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
	"log/slog"
	"net/http"
)

type Middleware struct {
	repo SessionRepository
	l    *slog.Logger
}

func NewMiddleware(r SessionRepository, l *slog.Logger) *Middleware {
	return &Middleware{
		repo: r,
		l:    l,
	}
}

type SessionRepository interface {
	GetSession(ctx context.Context, sessionID uuid.UUID) (entity.User, error)
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := m.l.With("request_id", uuid.NewString())
		l.Info("incoming request", "method", r.Method, "path", r.URL.Path)
		ctx := context.WithValue(r.Context(), "logger", l)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)

	})
}

func (m *Middleware) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("session_id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				HandlerError(ctx, w, entity.ErrNotAuthenticated)
				return
			}
			HandlerError(ctx, w, err)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value) //записанное значение Cookie при создании клиента записываем в sessionID
		if err != nil {
			HandlerError(ctx, w, err)
			return
		}

		user, err := m.repo.GetSession(ctx, sessionID) //в таблице связи  возвращаем нужный userID
		if err != nil {
			HandlerError(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, "user", user) //вносим в context
		r = r.WithContext(ctx)                     //перезаписываем запрос с новым контестом, в который сохранили userID
		next.ServeHTTP(w, r)
	})
}
