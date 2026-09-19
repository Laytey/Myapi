package service

import (
	"Myapi/internal/models"
	"Myapi/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

// 1. проверка на Email уже занят
func TestCreateUser_EmailExists(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("GetByEmail", "test@example.com").Return(models.User{Email: "test@example.com"}, nil)

	svc := NewUserService(mockRepo)

	_, err := svc.CreateUser(models.User{Email: "test@example.com"})
	if err == nil {
		t.Errorf("ожидали ошибку 'email already exists', получили nil")
	}
}

// 2. проверка имя уже занято
func TestCreateUser_NameExists(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("GetByEmail", "new@example.com").Return(models.User{}, errors.New("not found"))
	mockRepo.On("GetAll").Return([]models.User{{Name: "Existing"}}, nil)

	svc := NewUserService(mockRepo)

	_, err := svc.CreateUser(models.User{Email: "new@example.com", Name: "Existing"})
	if err == nil {
		t.Errorf("ожидали ошибку 'name already exists', получили nil")
	}
}

// 3. проверка занят пароль

func TestCreateUser_PasswordExists(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("GetByEmail", "new@example.com").Return(models.User{}, errors.New("not found"))
	mockRepo.On("GetAll").Return([]models.User{{Name: "Other", Password: "secret"}}, nil)

	svc := NewUserService(mockRepo)

	_, err := svc.CreateUser(models.User{
		Email:    "new@example.com",
		Name:     "NewName",
		Password: "secret",
	})
	if err == nil {
		t.Errorf("ожидали ошибку 'password already exists', получили nil")
	}
}

// 4. проверка на успешное создание пользователя

func TestCreateUser_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("GetByEmail", "new@example.com").Return(models.User{}, errors.New("not found"))
	mockRepo.On("GetAll").Return([]models.User{}, nil)
	mockRepo.On("Save", mock.Anything).Return(nil)

	svc := NewUserService(mockRepo)

	user, err := svc.CreateUser(models.User{
		Email:    "new@example.com",
		Name:     "NewName",
		Password: "secret",
	})
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if user.Email != "new@example.com" {
		t.Errorf("ожидали Email='new@example.com', получили %q", user.Email)
	}
}
