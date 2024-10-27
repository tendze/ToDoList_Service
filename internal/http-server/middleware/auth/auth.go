package auth

import (
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = r.Header.Get("Authorization")

		})
	}
}
