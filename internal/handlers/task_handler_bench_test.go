package handlers

import (
	"bytes"
	"net/http/httptest"
	"strconv"
	"testing"

	"Myapi/internal/models"
	"Myapi/internal/repository"
	"Myapi/internal/service"

	"github.com/gin-gonic/gin"
)

func BenchmarkCreateTask(b *testing.B) {

	// SETUP (один раз)
	gin.SetMode(gin.TestMode) // отключает debug-логи Gin
	// иначе каждая итерация будет печатать что-то в консоль - вывод засорится, замер исказится

	// создаем зависимости
	repo := repository.NewMemoryTaskRepository()
	svc := service.NewTaskServiceForTest(repo)
	handler := NewTaskHandler(svc)

	body := []byte(`{"title":"Test task","description":"Test description"}`)

	for b.Loop() {
		// b.N - Go сам подбирает число итераций
		req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		//ShouldBindJSON проверяет Content-Type

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w) // создать контекст Gin с фейковым writer
		c.Request = req                  // привязать запрос
		c.Set("userID", "user1")         // эмулировать работу AuthMiddleware. Без этого CreateTask вернёт 401

		handler.CreateTask(c)
	}

}

func BenchmarkGetTaskByID(b *testing.B) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryTaskRepository()
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	created, err := svc.CreateTask(models.Task{Title: "Test", UserUID: "user1"})
	if err != nil {
		b.Fatalf("не удалось создать задачу: %v", err)
	}

	idStr := strconv.Itoa(created.ID)

	for b.Loop() {
		req := httptest.NewRequest("GET", "/tasks/"+idStr, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: idStr}}
		c.Set("userID", "user1")

		handler.GetTaskByID(c)
	}
}

func BenchmarkLogin(b *testing.B) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryUserRepository()
	svc := service.NewUserService(repo)
	handler := NewAuthHandler(svc)

	// создаём пользователя в setup
	_, err := svc.CreateUser(models.User{
		Email:    "test@example.com",
		Password: "secret123",
	})
	if err != nil {
		b.Fatalf("не удалось создать пользователя: %v", err)
	}

	body := []byte(`{"email":"test@example.com","password":"secret123"}`)

	for b.Loop() {
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Login(c)
	}
}
