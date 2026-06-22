package jwt

import (
	"auth/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user model.User, app model.App, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid":    user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(duration).Unix(),
		"app_id": app.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(app.SecretKey)
}
