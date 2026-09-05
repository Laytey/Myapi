package service

import (
	"Myapi/internal/models"
	"Myapi/internal/repository"
	"errors"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task models.Task) (models.Task, error) {
	// Валидация: заголовок не может быть пустым
	if task.Title == "" {
		//чтобы не сохранять пустые задачи
		return models.Task{}, errors.New("title is required")
	}

	// Если статус не указан, ставим "Новая" по умолчанию
	if task.Status == "" {
		task.Status = "Новая"
	}

	err := s.repo.Save(&task)
	return task, err
}

func (s *TaskService) GetTaskByID(id int) (models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) GetTasksByUser(uid string) ([]models.Task, error) {
	return s.repo.GetByUserUID(uid)
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) UpdateTask(task models.Task) error {
	// Проверяем, существует ли задача
	_, err := s.repo.GetByID(task.ID)
	if err != nil {
		return errors.New("task not found")
	}

	return s.repo.Update(task)
}

func (s *TaskService) DeleteTask(id int) error {
	// Проверяем, существует ли задача
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	return s.repo.Delete(id)
}
