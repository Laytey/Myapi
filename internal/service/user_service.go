package service

import (
	"Myapi/internal/models"
	"Myapi/internal/repository"
	"errors"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

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

	err = s.repo.Save(user)
	return user, err
}

func (s *UserService) GetUserByID(id int) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) UpdateUser(user models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
