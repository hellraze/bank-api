package middleware

import (
	"fmt"
	"net/http"
)

func Recover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(writer, fmt.Sprintf("PANIC: %w", err), http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(writer, request)
	})
}
