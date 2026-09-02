package main

import (
	"Myapi/internal/auth"
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
	authHandler := handlers.NewAuthHandler(userService)

	r.POST("/login", authHandler.Login)
	r.POST("/users", userHandler.CreateUser) // регистрация без токена

	protected := r.Group("/") // создаем подроутер, сперва работает AuthMiddleware, потом Profile, GetAllUsers ...
	protected.Use(auth.AuthMiddleware)
	{
		protected.GET("/profile", authHandler.Profile)

		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// запуск сервера
	r.Run(":8080")
}
