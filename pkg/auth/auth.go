package auth

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenExp = 8 * time.Hour // 8 часов
)

type Claims struct {
	jwt.StandardClaims
}

func GenerateToken(password string) (string, error) {
	secret := os.Getenv("TODO_PASSWORD")
	if secret == "" {
		return "", fmt.Errorf("пароль не задан")
	}

	// Генерируем хеш из пароля
	hash := fmt.Sprintf("%x", hashPassword(password))

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExp).Unix(),
			Issuer:    "todo-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(hash))
}

func hashPassword(pass string) []byte {
	h := sha256.Sum256([]byte(pass))
	return h[:]
}
