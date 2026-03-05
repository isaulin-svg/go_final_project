package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Password string `json:"password"`
}

// GetPasswordHash возвращает SHA256 хеш пароля
func GetPasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// SigninHandler проверяет пароль и выдает JWT
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка JSON"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	// Создаем токен. В Claims кладем хеш пароля
	claims := jwt.MapClaims{
		"hash": GetPasswordHash(pass),
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(pass)) // Секрет - сам пароль
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка создания токена"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": signedToken})
}

// AuthMiddleware проверяет наличие и валидность токена
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		// Если пароль установлен в системе, проверяем авторизацию
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
				return
			}

			tokenStr := cookie.Value
			// Проверяем подпись токена
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(pass), nil
			})

			if err != nil || !token.Valid {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
				return
			}
		}

		// Если пароля нет в env или токен валиден — пропускаем дальше
		next(w, r)
	})
}
