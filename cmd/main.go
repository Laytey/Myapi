package main

import (
	"Myapi/internal/handlers"
	"Myapi/internal/repository"
	"Myapi/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// новый роутер
	r := gin.Default() // функция из gin, которая создаёт роутер с:
	// логированием (каждый запрос печатается в консоль)
	// восстановлением после паники?

	// Создаём зависимости
	userRepo := repository.NewMemoryUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.POST("/tasks", handlers.CreateTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	// Маршруты
	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.POST("/users", userHandler.CreateUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	// запуск сервера
	r.Run(":8080")
}
