package service

import (
	"errors"

	"Myapi/internal/models"
	"Myapi/internal/repository"
)

// UserService предоставляет операции над пользователями.
//
// Работает поверх UserRepository. Не зависит от HTTP-слоя.
// Потокобезопасен: безопасность обеспечивается реализацией репозитория.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService создаёт сервис пользователей.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser создаёт нового пользователя.
//
// Проверяет уникальность email, name и password:
//   - Если email занят — возвращает "email already exists".
//   - Если name занят — возвращает "name already exists".
//   - Если такой пароль уже есть — возвращает "password already exists".
//
// Возвращает созданного пользователя или ошибку.
//
// NOTE: проверка уникальности пароля — учебный вариант решения.
// В продакшене пароли хешируются (bcrypt/argon2), и одинаковые
// пароли дают разные хеши. Эту проверку нужно убрать.
func (s *UserService) CreateUser(user models.User) (models.User, error) {

	_, err := s.repo.GetByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("email already exists")
	}

	users, _ := s.repo.GetAll()
	for _, u := range users {
		if u.Name == user.Name {
			return models.User{}, errors.New("name already exists")
		}
	}

	for _, u := range users {
		if u.Password == user.Password {
			return models.User{}, errors.New("password already exists")
		}
	}

	err = s.repo.Save(&user)
	return user, err
}

// GetUserByID возвращает пользователя по его ID.
func (s *UserService) GetUserByID(id int) (models.User, error) {
	return s.repo.GetByID(id)
}

// GetAllUsers возвращает всех пользователей.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

// UpdateUser обновляет существующего пользователя.
func (s *UserService) UpdateUser(user models.User) error {
	return s.repo.Update(user)
}

// DeleteUser удаляет пользователя по ID.
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// GetUserByEmail возвращает пользователя по email.
// Используется при логине и проверке уникальности при регистрации.
func (s *UserService) GetUserByEmail(email string) (models.User, error) {
	return s.repo.GetByEmail(email)
}
