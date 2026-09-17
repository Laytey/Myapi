package main

import (
	"Myapi/internal/auth"
	"Myapi/internal/db"
	"Myapi/internal/handlers"
	"Myapi/internal/middleware"
	"Myapi/internal/repository"
	"Myapi/internal/service"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// новый роутер
	r := gin.Default()
	// gin.Default() — это стандартный HTTP-маршрутизатор, который предоставляет базовые функции для обработки HTTP
	// функция из gin, которая создаёт роутер с двумя middleware
	// логированием (каждый запрос печатается в консоль)
	// восстановлением после паники?
	// попытка подключиться к PostgreSQL
	connString := os.Getenv("DATABASE_URL")
	// os - стандартный пакет Go для работы
	// дает доступ к переменным окружения (Getenv), аргументам командной строки, файловой системе
	if connString == "" {
		connString = "postgres://user:password@localhost:5432/mydb?sslmode=disable"
	}

	storage, err := db.New(context.Background(), connString)

	var userRepo repository.UserRepository

	var taskRepo repository.TaskRepository

	if err != nil {
		log.Println("PostgreSQL not available, using memory storage")
		userRepo = repository.NewMemoryUserRepository()
		taskRepo = repository.NewMemoryTaskRepository()
	} else {
		log.Println("Using PostgreSQL storage")
		userRepo = repository.NewPostgresUserRepository(storage)
		taskRepo = repository.NewPostgresTaskRepository(storage)
	}

	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	// Создаём зависимости

	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)
	authHandler := handlers.NewAuthHandler(userService)

	r.Use(middleware.GzipMiddleware)

	r.POST("/login", authHandler.Login)      // не защищаем, чтобы можно было зайти без токена и его получить
	r.POST("/users", userHandler.CreateUser) // регистрация тоже без токена

	protected := r.Group("/") // создаем подроутер, сперва работает AuthMiddleware, потом Profile, GetAllUsers ...
	protected.Use(auth.AuthMiddleware)
	{
		protected.GET("/profile", authHandler.Profile)

		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		protected.GET("/tasks", taskHandler.GetAllTasks)
		protected.GET("/tasks/:id", taskHandler.GetTaskByID)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
	}

	// запуск сервера
	r.Run(":8080")
}
