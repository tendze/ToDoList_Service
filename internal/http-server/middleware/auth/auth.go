package auth

import (
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
	"todolist/internal/client/auth"
	"todolist/internal/lib/api/response"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	log = log.With(
		slog.String("component", "middleware/auth"),
	)
	log.Info("api protection middleware enabled")
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("no token provided")
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("no token provided"))
				return
			}
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				log.Error("Invalid authorization header format")
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid authorization header format"))
				return
			}
			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			resp, err := auth.ValidateToken(tokenString)
			if err != nil {
				log.Error("failed to validate token", err)
				w.WriteHeader(http.StatusBadGateway)
				render.JSON(w, r, response.Error("failed to validate token"))
				return
			}
			ctx := context.WithValue(r.Context(), "login", resp.UserLogin)
			log.Info("token validated. user login:", resp.UserLogin)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
