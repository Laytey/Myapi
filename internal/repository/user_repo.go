package repository

import "Myapi/internal/models"

// UserRepository — интерфейс для работы с пользователями
type UserRepository interface {
	Save(user models.User) error
	GetByID(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetAll() ([]models.User, error)
	Update(user models.User) error
	Delete(id int) error
}
