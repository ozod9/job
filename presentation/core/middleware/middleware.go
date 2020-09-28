package middleware

import (
	"net/http"
)

func Requests(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := StatusRecorder{ResponseWriter: w, Status: http.StatusOK}
		next.ServeHTTP(&sw, r)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.Status = code
	rec.ResponseWriter.WriteHeader(code)
}
