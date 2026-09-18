package service

import (
	"Myapi/internal/models"
	"Myapi/internal/repository"
	"errors"
	"log"
	"time"
)

type TaskService struct {
	repo      repository.TaskRepository
	cleanupCh chan struct{}
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	s := &TaskService{
		repo:      repo,
		cleanupCh: make(chan struct{}, 10),
		//для накопления 10 сигналов
		//если буфер заполнен — значит, накопилось 10 событий
	}
	go s.cleanupWorker()
	return s
}

func (s *TaskService) cleanupWorker() {
	for range s.cleanupCh { //цикл, который работает бесконечно, пока канал не закрыт
		time.Sleep(2 * time.Second)
		if len(s.cleanupCh) >= cap(s.cleanupCh)-1 {
			log.Println("Канал заполнен, запускаем hard delete")
			if err := s.repo.HardDelete(); err != nil {
				//объявление + проверка в одной строке. err существует только внутри if/else
				log.Println("Ошибка hard delete:", err)
			} else {
				log.Println("Hard delete выполнен успешно")
			}
			for len(s.cleanupCh) > 0 {
				<-s.cleanupCh
			}
		}

	}
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
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.cleanupCh <- struct{}{}
	return nil
}
