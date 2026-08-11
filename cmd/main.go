package main

import (
	"Myapi/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// новый роутер
	r := gin.Default() // функция из gin, которая создаёт роутер с:
	// логированием (каждый запрос печатается в консоль)
	// восстановлением после паники?

	// включаем маршруты (роутеры)
	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/tasks/:id", handlers.GetTaskByID)
	r.POST("/tasks", handlers.CreateTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.GET("/users", handlers.GetAllUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	// запуск сервера
	r.Run(":8080")
}
