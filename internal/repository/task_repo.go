// Package repository определяет интерфейсы для доступа к данным
// и их реализации: in-memory (для тестов и разработки без БД)
// и PostgreSQL (для продакшена).
package repository

import "Myapi/internal/models"

// TaskRepository определяет контракт для работы с задачами.
//
// Все методы должны быть безопасны при конкурентном доступе,
// если реализация используется из нескольких горутин.
type TaskRepository interface {
	// Save сохраняет новую задачу и присваивает ей уникальный ID.
	// ID записывается в переданную структуру через указатель.
	Save(task *models.Task) error

	// GetByID возвращает задачу по её уникальному ID.
	// Возвращает ошибку "task not found", если задача не существует
	// или помечена как удалённая.
	GetByID(id int) (models.Task, error)

	// GetByUserUID возвращает все активные задачи указанного пользователя.
	// Удалённые задачи (Deleted = true) не включаются в результат.
	GetByUserUID(uid string) ([]models.Task, error)

	// GetAll возвращает все активные задачи всех пользователей.
	// Удалённые задачи (Deleted = true) не включаются в результат.
	GetAll() ([]models.Task, error)

	// Update обновляет существующую задачу.
	// Возвращает ошибку, если задача с таким ID не найдена.
	Update(task models.Task) error

	// Delete помечает задачу как удалённую (soft-delete).
	// Физического удаления не происходит — задача остаётся в БД
	// с флагом Deleted = true.
	Delete(id int) error

	// HardDelete физически удаляет все задачи с Deleted = true.
	// Используется для периодической очистки в фоне.
	HardDelete() error
}
