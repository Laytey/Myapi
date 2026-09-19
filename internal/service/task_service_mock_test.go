package service

import (
	"errors"
	"testing"
	"time"

	"Myapi/internal/models"
	"Myapi/internal/repository/mocks"

	"github.com/stretchr/testify/mock"
)

func TestCreateTask_SaveError(t *testing.T) {
	// 1. Создаём мок-репозиторий
	mockRepo := mocks.NewTaskRepository(t)

	// 2. Настраиваем поведение: при вызове Save вернуть ошибку
	mockRepo.On("Save", mock.Anything).Return(errors.New("db error"))

	// 3. Создаём сервис с моком
	svc := NewTaskServiceForTest(mockRepo)

	// 4. Вызываем CreateTask
	_, err := svc.CreateTask(models.Task{Title: "Test"})

	// 5. Проверяем, что ошибка прокинулась
	if err == nil {
		t.Errorf("ожидали ошибку от Save, получили nil")
	}
}

func TestDeleteTask_DeleteError(t *testing.T) {
	mockRepo := mocks.NewTaskRepository(t)

	// GetByID - вернёт задачу (она "существует")
	mockRepo.On("GetByID", 1).Return(models.Task{ID: 1, Title: "Test"}, nil)
	// Delete - вернёт ошибку
	mockRepo.On("Delete", 1).Return(errors.New("db error"))

	svc := NewTaskServiceForTest(mockRepo)

	err := svc.DeleteTask(1)
	if err == nil {
		t.Errorf("ожидали ошибку от Delete, получили nil")
	}
}

func TestCleanupWorker(t *testing.T) {
	mockRepo := mocks.NewTaskRepository(t)

	// HardDelete вызывается воркером при заполнении канала
	mockRepo.On("HardDelete").Return(errors.New("hard delete error"))

	// Создаём сервис С воркером (не ForTest!)
	svc := NewTaskService(mockRepo)

	// Записываем 10 сигналов напрямую в канал (обходя DeleteTask)
	for i := 0; i < 10; i++ {
		svc.cleanupCh <- struct{}{}
	}

	// Ждём, пока воркер проснётся (Sleep 2s) + запас
	time.Sleep(3 * time.Second)

	// Проверяем, что HardDelete был вызван
	mockRepo.AssertCalled(t, "HardDelete")
}
func TestNewTaskService(t *testing.T) {
	mockRepo := mocks.NewTaskRepository(t)
	svc := NewTaskService(mockRepo)

	if svc == nil {
		t.Fatal("NewTaskService вернул nil")
	}
	if svc.cleanupCh == nil {
		t.Error("cleanupCh не инициализирован")
	}
}

func TestGetTaskByID_Error(t *testing.T) {
	mockRepo := mocks.NewTaskRepository(t)
	mockRepo.On("GetByID", 1).Return(models.Task{}, errors.New("db error"))

	svc := NewTaskServiceForTest(mockRepo)

	_, err := svc.GetTaskByID(1)
	if err == nil {
		t.Errorf("ожидали ошибку от GetByID, получили nil")
	}
}

func TestGetTasksByUser_Error(t *testing.T) {
	mockRepo := mocks.NewTaskRepository(t)
	mockRepo.On("GetByUserUID", "user1").Return(nil, errors.New("db error"))

	svc := NewTaskServiceForTest(mockRepo)

	_, err := svc.GetTasksByUser("user1")
	if err == nil {
		t.Errorf("ожидали ошибку от GetByUserUID, получили nil")
	}
}

func TestGetAllTasks_Error(t *testing.T) {
	mockRepo := mocks.NewTaskRepository(t)
	mockRepo.On("GetAll").Return(nil, errors.New("db error"))

	svc := NewTaskServiceForTest(mockRepo)

	_, err := svc.GetAllTasks()
	if err == nil {
		t.Errorf("ожидали ошибку от GetAll, получили nil")
	}
}
