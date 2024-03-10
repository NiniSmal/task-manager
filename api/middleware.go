package api

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/repository"
	"log"
	"net/http"
)

type Middleware struct {
	repo *repository.TaskRepository
}

func NewMiddleware(r *repository.TaskRepository) *Middleware {
	return &Middleware{
		repo: r,
	}
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
			HandlerError(w, err)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value) //записанное значение Cookie при создании клиента записываем в sessionID
		if err != nil {
			HandlerError(w, err)
			return
		}

		userID, role, err := m.repo.GetUserIDBySessionID(ctx, sessionID) //в таблице связи  возвращаем нужный userID
		if err != nil {
			HandlerError(w, err)
			return
		}

		ctx = context.WithValue(ctx, "user_id", userID) //вносим в context
		ctx = context.WithValue(ctx, "role", role)
		r = r.WithContext(ctx) //перезаписываем запрос с новым контестом, в который сохранили userID
		next.ServeHTTP(w, r)
	})
}

//добавляет заголовок к ответу.

//type ResponseHeader struct {
//	handler     http.Handler
//	headerName  string
//	headerValue string
//}

// обработчик
//func NewResponseHeader(handlerToWrap http.Handler, headerName string, headerValue string) *ResponseHeader {
//	return &ResponseHeader{handlerToWrap, headerName, headerValue}
//}
