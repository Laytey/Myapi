package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	// способ вытащить значение заголовка Authorization из HTTP-запроса
	// там должен быть токен
	if authHeader == "" {
		ctx.AbortWithStatusJSON( // AbortWithStatusJSON прерывает выполнение запроса
			http.StatusUnauthorized,
			gin.H{
				"error": "authorization header is required",
			},
		)
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" { // Bearer — стандартный префикс для JWT токенов
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid authorisation header",
			},
		)
		return
	}
	claims, err := ParseToken(parts[1])
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
