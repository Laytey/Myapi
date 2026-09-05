package handlers

import (
	"Myapi/internal/models"
	"Myapi/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// возвращаем все задачи - клиент запрашивает данные
// Теперь: только задачи текущего пользователя
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	// Получаем UID текущего пользователя из контекста
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	// Получаем задачи пользователя через сервис
	tasks, err := h.service.GetTasksByUser(uidStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GET /tasks/:id — получить задачу по ID
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id") // извлекает параметр id из URL?
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Проверяем, что задача существует
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, что задача принадлежит текущему пользователю
	uid, _ := c.Get("userID")
	uidStr, _ := uid.(string)

	if task.UserUID != uidStr {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// создаём новую задачу - POST
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task

	// читаем JSON из запроса - отправляет в переменную task?
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Привязываем задачу к текущему пользователю
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	task.UserUID = uidStr

	// создаем ID - присваиваем значения полю ID структуры task?...
	// (теперь это делает репозиторий, но оставляем комментарий)

	// добавляем в хранилище через сервис
	createdTask, err := h.service.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// возвращаем созданную задачу - присылаем ответ в формате JSON и статус 201
	c.JSON(http.StatusCreated, createdTask)
}

// PUT /tasks/:id — обновить задачу
func (h *TaskHandler) UpdateTask(c *gin.Context) {
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

	// Проверяем, что задача существует и принадлежит пользователю
	existingTask, err := h.service.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	uid, _ := c.Get("userID")
	uidStr, _ := uid.(string)

	if existingTask.UserUID != uidStr {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Обновляем поля, которые передал клиент
	if updatedTask.Title != "" {
		existingTask.Title = updatedTask.Title
	}
	if updatedTask.Description != "" {
		existingTask.Description = updatedTask.Description
	}
	if updatedTask.Status != "" {
		existingTask.Status = updatedTask.Status
	}

	err = h.service.UpdateTask(existingTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingTask)
}

// DELETE /tasks/:id — удалить задачу
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Проверяем, что задача существует и принадлежит пользователю
	existingTask, err := h.service.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	uid, _ := c.Get("userID")
	uidStr, _ := uid.(string)

	if existingTask.UserUID != uidStr {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204 no content - успешно, но без отправки данных
}
