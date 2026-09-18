package repository

import (
	"Myapi/internal/models"
	"errors"
)

type MemoryTaskRepository struct {
	tasks  []models.Task
	nextID int
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks:  []models.Task{},
		nextID: 1, // начальное значение ID для новой задачи
	}
}

func (r *MemoryTaskRepository) Save(task *models.Task) error {
	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, *task)
	return nil
}

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

func (r *MemoryTaskRepository) GetByUserUID(uid string) ([]models.Task, error) {
	var result []models.Task
	for _, t := range r.tasks {
		if t.UserUID == uid && !t.Deleted {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *MemoryTaskRepository) GetAll() ([]models.Task, error) {
	var result []models.Task
	for _, t := range r.tasks {
		if !t.Deleted {
			result = append(result, t)
		}
	}
	return result, nil

}

func (r *MemoryTaskRepository) Update(task models.Task) error {
	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = task
			return nil
		}
	}
	return errors.New("task not found")
}

func (r *MemoryTaskRepository) Delete(id int) error {
	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks[i].Deleted = true
			return nil
		}
	}
	return errors.New("task not found")
}

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
