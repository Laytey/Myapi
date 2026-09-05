package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	var token string

	// Проверяем заголовок Authorization
	authHeader := ctx.GetHeader("Authorization")
	// способ вытащить значение заголовка Authorization из HTTP-запроса
	// там должен быть токен
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			// Bearer — стандартный префикс для JWT токенов
			token = parts[1]
		}
	}

	// Если в заголовке нет — проверяем cookie
	if token == "" {
		cookie, err := ctx.Cookie("token")
		if err == nil {
			token = cookie
		}
	}

	// Если токена нет нигде — запрет доступа
	if token == "" {
		ctx.AbortWithStatusJSON( // AbortWithStatusJSON прерывает выполнение запроса
			http.StatusUnauthorized,
			gin.H{
				"error": "authorization required",
			},
		)
		return
	}

	// Проверяем токен
	claims, err := ParseToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.Set("userID", claims.UserID) // Set сохраняет значение в контексте запроса
	ctx.Next()                       // передаёт управление следующему обработчику
}
