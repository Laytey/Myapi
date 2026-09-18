package repository

import "Myapi/internal/models"

type TaskRepository interface {
	Save(task *models.Task) error
	GetByID(id int) (models.Task, error)
	GetByUserUID(uid string) ([]models.Task, error) // задачи конкретного пользователя
	GetAll() ([]models.Task, error)
	Update(task models.Task) error
	Delete(id int) error
	HardDelete() error
}
