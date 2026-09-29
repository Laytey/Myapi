package handlers

import (
	"net/http"
	"strconv"
	"time"

	"Myapi/internal/auth"
	"Myapi/internal/models"
	"Myapi/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает HTTP-запросы аутентификации и профиля.
//
// Делегирует работу с пользователями в UserService. Не содержит
// бизнес-логики - только координацию: парсинг, вызов сервиса,
// формирование ответа.
//
// NOTE: пароли в этой реализации сравниваются в открытом виде
// (в UserService.CreateUser). В рабочих проектах требуется хеширование. Обратить внимание.
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler создаёт хендлер аутентификации с указанным сервисом.
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Login обрабатывает POST /login.
// Проверяет email и пароль, при успехе генерирует JWT и устанавливает
// его как cookie "token" (HttpOnly).
//
// Тело запроса — JSON с полями email, password.
//
// Ответы:
//   - 400 Bad Request — если JSON некорректен.
//   - 401 Unauthorized — если email или пароль не совпадают.
//   - 500 Internal Server Error — при ошибке генерации токена.
//   - 200 OK — {"token": "<jwt>"} + cookie "token".
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(strconv.Itoa(user.ID), 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 86400, "/", "", false, true)
	// метод Gin для установки cookie в браузере клиента
	// token-Сам JWT-токен (строка)
	// Время жизни cookie в секундах 86400 (1 день)
	//"/" домен не указан: действует для текущего хоста (localhost)
	// "" домен, для которого действует cookie (пусто = текущий домен)
	// false	Не требуем HTTPS (для локальной разработки)
	// true	Защита: JavaScript не сможет прочитать этот cookie

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Profile обрабатывает GET /profile.
// Возвращает профиль текущего аутентифицированного пользователя.
//
// Требует, чтобы AuthMiddleware установил userID в контексте Gin.
//
// Возвращает только публичные поля (id, name, email) — без password.
//
// Ответы:
//   - 401 Unauthorized — если userID отсутствует.
//   - 400 Bad Request — если userID не число.
//   - 404 Not Found — если пользователь не найден.
//   - 200 OK — {"id": ..., "name": ..., "email": ...}.
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userService.GetUserByID(idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
