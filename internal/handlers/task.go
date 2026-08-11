package handlers

import (
	"Myapi/internal/models"
	"Myapi/internal/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// возвращаем все задачи - клиент запрашивает данные
func GetAllTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, storage.Tasks)
}

// создаём новую задачу - POST
func CreateTask(ctx *gin.Context) {
	var task models.Task

	// читаем JSON из запроса - отправляет в переменную tаsk?
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем ID - присваиваем значения полю ID структуры task?...
	task.ID = len(storage.Tasks) + 1

	// добавляем в хранилище
	storage.Tasks = append(storage.Tasks, task)

	// возвращаем созданную задачу - присылаем ответ в формате JSON и статус 201
	ctx.JSON(http.StatusCreated, task)
}
func GetTaskByID(c *gin.Context) {
	idStr := c.Param("id") // извлекает параметр id из URL?
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, task := range storage.Tasks {
		if task.ID == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range storage.Tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				storage.Tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				storage.Tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				storage.Tasks[i].Status = updatedTask.Status
			}
			c.JSON(http.StatusOK, storage.Tasks[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, task := range storage.Tasks {
		if task.ID == id {
			storage.Tasks = append(storage.Tasks[:i], storage.Tasks[i+1:]...)
			// ... — это оператор распаковки, чтобы превратить слайс в отдельные элементы
			// берём все элементы до индекса и после него и (...) склеиваем эти две части
			// получится слайс без того элемента, кторый надо удалить
			c.JSON(http.StatusNoContent, nil) // 204 no content - успешно, но без отправки данных
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}
