package middleware

import (
	"log"
	"net/http"
)

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
