package handlers

import (
	"net/http"
	"strconv"

	"Myapi/internal/models"
	"Myapi/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler обрабатывает HTTP-запросы, связанные с пользователями.
//
// Делегирует (передает) бизнес-логику в UserService.
//
// NOTE: в текущей реализации не проверяется, что запрашиваемый
// пользователь совпадает с аутентифицированным. Это означает, что
// любой авторизованный пользователь может обновлять и удалять
// других пользователей. В реальном проекте требуется проверка прав.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler создаёт хендлер пользователей с указанным сервисом.
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUsers обрабатывает GET /users.
// Возвращает список всех пользователей.
//
// Ответы:
//   - 500 Internal Server Error — при ошибке сервиса.
//   - 200 OK — список пользователей в JSON.
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID обрабатывает GET /users/:id.
// Возвращает пользователя по ID.
//
// Ответы:
//   - 400 Bad Request — если ID не число.
//   - 404 Not Found — если пользователь не найден.
//   - 200 OK — пользователь в JSON.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser обрабатывает POST /users.
// Регистрирует нового пользователя. Не требует аутентификации.
//
// Тело запроса — JSON с полями name, email, password.
//
// Ответы:
//   - 400 Bad Request — если JSON некорректен.
//   - 409 Conflict — если email, name или password уже заняты.
//   - 201 Created — созданный пользователь в JSON.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

// UpdateUser обрабатывает PUT /users/:id.
// Обновляет пользователя. ID из URL имеет приоритет над ID в теле.
//
// Ответы:
//   - 400 Bad Request — если ID или JSON некорректны.
//   - 404 Not Found — если пользователь не найден.
//   - 200 OK — обновлённый пользователь в JSON.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id

	err = h.service.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser обрабатывает DELETE /users/:id.
// Физически удаляет пользователя из хранилища.
//
// Ответы:
//   - 400 Bad Request — если ID не число.
//   - 404 Not Found — если пользователь не найден.
//   - 204 No Content — успешное удаление.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
