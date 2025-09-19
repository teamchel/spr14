package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Проверяем токен
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			hash := hashPassword(pass)
			return hash, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
