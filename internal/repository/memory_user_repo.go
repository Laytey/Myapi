package repository

import (
	"errors"

	"Myapi/internal/models"
)

// MemoryUserRepository — реализация UserRepository, хранящая
// пользователей в срезе в оперативной памяти.
//
// Подходит для тестов и разработки без БД. Данные теряются при
// перезапуске процесса. Не безопасна при конкурентном доступе
// без внешней синхронизации.
type MemoryUserRepository struct {
	users  []models.User
	nextID int
}

// NewMemoryUserRepository создаёт пустой in-memory репозиторий
// с начальным значением nextID = 1.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{ //
		users:  []models.User{}, // пустой слайс
		nextID: 1,
	}
}

// Save реализует UserRepository.Save.
// Присваивает пользователю следующий свободный ID и добавляет его в срез.
func (r *MemoryUserRepository) Save(user *models.User) error {
	user.ID = r.nextID
	r.nextID++
	r.users = append(r.users, *user)
	return nil // возвращает error.nil
}

// GetByID реализует UserRepository.GetByID.
// Возвращает ошибку "user not found", если пользователь не найден.
func (r *MemoryUserRepository) GetByID(id int) (models.User, error) {
	for _, u := range r.users { // u -каждый пользователь в цикле
		if u.ID == id {
			return u, nil
		}
	} // models.User{} - пустой пользователь
	return models.User{}, errors.New("user not found")
}

// GetByEmail реализует UserRepository.GetByEmail.
// Возвращает ошибку "user not found", если пользователь с таким email
// не зарегистрирован.
func (r *MemoryUserRepository) GetByEmail(email string) (models.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

// GetAll реализует UserRepository.GetAll.
func (r *MemoryUserRepository) GetAll() ([]models.User, error) {
	return r.users, nil
}

// Update реализует UserRepository.Update.
// Перезаписывает пользователя с тем же ID целиком.
func (r *MemoryUserRepository) Update(user models.User) error {
	for i, u := range r.users { // user.ID - ID пользователя, кот мы обновляем
		if u.ID == user.ID { //u.ID - ID пользователя из хранилища
			r.users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}

// Delete реализует UserRepository.Delete.
// В отличие от MemoryTaskRepository, выполняет физическое удаление:
// пользователь полностью убирается из среза.
func (r *MemoryUserRepository) Delete(id int) error {
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...) // - физическое удаление
			return nil
		}
	}
	return errors.New("user not found")
}
