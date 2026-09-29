// Package service содержит бизнес-логику приложения: операции над
// задачами и пользователями. Сервисы работают поверх репозиториев
// и не зависят от HTTP-слоя.
package service

import (
	"errors"
	"log"
	"time"

	"Myapi/internal/models"
	"Myapi/internal/repository"
)

// TaskService предоставляет операции над задачами.
//
// Работает поверх TaskRepository. Дополнительно управляет фоновым
// воркером очистки, который при накоплении сигналов вызывает
// физическое удаление soft-deleted задач.
//
// Потокобезопасен: методы могут вызываться из нескольких горутин
// одновременно (безопасность обеспечивается реализацией репозитория).
type TaskService struct {
	repo      repository.TaskRepository
	cleanupCh chan struct{}
}

// NewTaskService создаёт сервис задач и запускает фоновый воркер
// очистки. Воркер реагирует на сигналы в cleanupCh: когда в канале
// накопилось до 10 сигналов - вызывает repo.HardDelete.
//
// Воркер работает до конца жизни приложения. Для тестов используй
// NewTaskServiceForTest - он не запускает воркер.
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
				<-s.cleanupCh // вычитываем из канала - очистка
			}
		}

	}
}

// CreateTask создаёт новую задачу.
//
// Валидация:
//   - Title не может быть пустым.
//   - Если Status не указан - по умолчанию "Новая".
//
// Возвращает созданную задачу с присвоенным ID или ошибку.
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

// GetTaskByID возвращает задачу по её ID.
// Проксирует вызов в repo.GetByID.
func (s *TaskService) GetTaskByID(id int) (models.Task, error) {
	return s.repo.GetByID(id)
}

// GetTasksByUser возвращает все активные задачи указанного пользователя.
func (s *TaskService) GetTasksByUser(uid string) ([]models.Task, error) {
	return s.repo.GetByUserUID(uid)
}

// GetAllTasks возвращает все активные задачи всех пользователей.
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAll()
}

// UpdateTask обновляет существующую задачу.
// Перед обновлением проверяет, что задача существует.
// Возвращает ошибку "task not found", если задачи с таким ID нет.
func (s *TaskService) UpdateTask(task models.Task) error {
	// Проверяем, существует ли задача
	_, err := s.repo.GetByID(task.ID)
	if err != nil {
		return errors.New("task not found")
	}

	return s.repo.Update(task)
}

// DeleteTask выполняет soft-delete задачи: помечает её удалённой
// и отправляет сигнал в cleanupCh для последующей физической очистки.
//
// Возвращает ошибку "task not found", если задачи не существует.
func (s *TaskService) DeleteTask(id int) error {
	// Проверяем, существует ли задача
	// не важно, что именно вернул GetByID - важно, что он вообще не должен найти задачу
	_, err := s.repo.GetByID(id)
	if err != nil { // err == nil — если ошибки нет, значит нашёл
		return errors.New("task not found")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.cleanupCh <- struct{}{}
	return nil
}

// NewTaskServiceForTest создаёт сервис без запуска фонового воркера.
// Используется в unit-тестах, чтобы не было фоновых горутин
// и можно было проверить содержимое cleanupCh.
func NewTaskServiceForTest(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
		// канал ссылочный тип -изменения в нём видны всем, кто держит копию структуры
		cleanupCh: make(chan struct{}, 10),
	}
}
