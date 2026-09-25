package handlers

import (
	"Myapi/internal/models"
	"Myapi/internal/repository"
	"Myapi/internal/service"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkGetAllTasks(b *testing.B) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryTaskRepository()
	svc := service.NewTaskServiceForTest(repo)
	handler := NewTaskHandler(svc)

	// создаём 10 задач в setup
	for i := 0; i < 10; i++ {
		svc.CreateTask(models.Task{Title: "Task", UserUID: "user1"})
	}

	for b.Loop() {
		req := httptest.NewRequest("GET", "/tasks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", "user1")

		handler.GetAllTasks(c)
	}
}
