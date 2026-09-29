package repository

import (
	"errors"

	"Myapi/internal/models"
)

// MemoryTaskRepository — реализация TaskRepository, хранящая задачи
// в срезе в оперативной памяти.
//
// Подходит для тестов и разработки без БД. Данные теряются при
// перезапуске процесса. Не безопасна при конкурентном доступе
// без внешней синхронизации.
type MemoryTaskRepository struct {
	tasks  []models.Task
	nextID int
}

// NewMemoryTaskRepository создаёт пустой in-memory репозиторий
// с начальным значением nextID = 1.
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks:  []models.Task{},
		nextID: 1, // начальное значение ID для новой задачи
	}
}

// Save реализует TaskRepository.Save.
// Присваивает задаче следующий свободный ID и добавляет её в срез.
func (r *MemoryTaskRepository) Save(task *models.Task) error {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, *task)
	return nil
}

// GetByID реализует TaskRepository.GetByID.
// Возвращает ошибку "task not found", если задача не найдена
// или помечена удалённой.
func (r *MemoryTaskRepository) GetByID(id int) (models.Task, error) {
	//ищет задачу по уникальному ID, и такой ID может быть только у одной задачи
	for _, t := range r.tasks {
		//t — это копия текущей задачи из слайса
		if t.ID == id && !t.Deleted {
			return t, nil
		}
	}
	return models.Task{}, errors.New("task not found")
}

// GetByUserUID реализует TaskRepository.GetByUserUID.
func (r *MemoryTaskRepository) GetByUserUID(uid string) ([]models.Task, error) {
	var result []models.Task
	for _, t := range r.tasks {
		if t.UserUID == uid && !t.Deleted {
			result = append(result, t)
		}
	}
	return result, nil
}

// GetAll реализует TaskRepository.GetAll.
func (r *MemoryTaskRepository) GetAll() ([]models.Task, error) {
	var result []models.Task
	for _, t := range r.tasks {
		if !t.Deleted {
			result = append(result, t)
		}
	}
	return result, nil

}

// Update реализует TaskRepository.Update.
// Перезаписывает задачу с тем же ID целиком.
func (r *MemoryTaskRepository) Update(task models.Task) error {
	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = task
			return nil
		}
	}
	return errors.New("task not found")
}

// Delete реализует TaskRepository.Delete.
// Помечает задачу как удалённую (soft-delete), не удаляя из среза.
func (r *MemoryTaskRepository) Delete(id int) error {
	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks[i].Deleted = true
			return nil
		}
	}
	return errors.New("task not found")
}

// HardDelete реализует TaskRepository.HardDelete.
// Физически удаляет все задачи с флагом Deleted = true,
// создавая новый срез только из «живых» задач.
func (r *MemoryTaskRepository) HardDelete() error {
	var alive []models.Task
	for _, t := range r.tasks {
		if !t.Deleted {
			alive = append(alive, t)
		}
	}
	r.tasks = alive
	return nil
}
