package handlers

import (
	"context"
	"net/http"

	"NazirovJasur/my-http-server.git/internal/repository"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(authRepo *repository.AuthRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userID, err := authRepo.GetUserBySession(cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Кладем ID пользователя в контекст запроса
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
