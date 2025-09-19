package api

import (
	"encoding/json"
	"net/http"
	"os"

	"spr14/pkg/auth"
)

type LoginReq struct {
	Password string `json:"password"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка парсинга JSON"})
		return
	}

	// Получаем сохранённый пароль из переменной окружения
	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		writeJSON(w, map[string]string{"error": "сервер не настроен"})
		return
	}

	if req.Password != expectedPass {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}

	// Генерируем JWT-токен
	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка генерации токена"})
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   int(8 * 3600), // 8 часов
		SameSite: http.SameSiteLaxMode,
	})

	// Ответ с токеном
	writeJSON(w, map[string]string{"token": token})
}
