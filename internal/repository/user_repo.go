package repository

import "Myapi/internal/models"

// UserRepository определяет контракт для работы с пользователями.
//
// Реализации: MemoryUserRepository (в памяти, для тестов)
// и PostgresUserRepository (PostgreSQL, для продакшена).
type UserRepository interface {
	// Save сохраняет нового пользователя и присваивает ему уникальный ID.
	// ID записывается в переданную структуру через указатель.
	Save(user *models.User) error

	// GetByID возвращает пользователя по уникальному ID.
	// Возвращает ошибку "user not found", если пользователь не существует.
	GetByID(id int) (models.User, error)

	// GetByEmail возвращает пользователя по email.
	// Используется при логине и проверке уникальности email при регистрации.
	GetByEmail(email string) (models.User, error)

	// GetAll возвращает всех пользователей.
	GetAll() ([]models.User, error)

	// Update обновляет существующего пользователя.
	// Возвращает ошибку, если пользователь с таким ID не найден.
	Update(user models.User) error

	// Delete удаляет пользователя по ID.
	// Возвращает ошибку, если пользователь не найден.
	Delete(id int) error
}
