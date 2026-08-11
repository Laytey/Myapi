package service

import (
	"Myapi/internal/models"
	"Myapi/internal/storage"
	"errors"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user models.User) (models.User, error) {

	for _, u := range storage.Users {
		if u.Email == user.Email {
			return models.User{}, errors.New("email already exists")
		}
	}

	for _, u := range storage.Users {
		if u.Name == user.Name {
			return models.User{}, errors.New("name already exists")
		}
	}

	for _, u := range storage.Users {
		if u.Password == user.Password {
			return models.User{}, errors.New("password already exists")
		}
	}

	user.ID = len(storage.Users) + 1
	storage.Users = append(storage.Users, user)
	return user, nil
}

func (s *UserService) GetAllUsers() []models.User {
	return storage.Users
}

func (s *UserService) GetUserByID(id int) (models.User, error) {
	for _, user := range storage.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (s *UserService) UpdateUser(id int, updatedUser models.User) (models.User, error) {
	for i, user := range storage.Users {
		if user.ID == id {

			if updatedUser.Name != "" {
				for _, u := range storage.Users {
					if u.Name == updatedUser.Name && u.ID != id {
						return models.User{}, errors.New("name already exists")
					}
				}
				storage.Users[i].Name = updatedUser.Name
			}

			if updatedUser.Email != "" {
				for _, u := range storage.Users {
					if u.Email == updatedUser.Email && u.ID != id {
						return models.User{}, errors.New("email already exists")
					}
				}
				storage.Users[i].Email = updatedUser.Email
			}

			if updatedUser.Password != "" {
				for _, u := range storage.Users {
					if u.Password == updatedUser.Password && u.ID != id {
						return models.User{}, errors.New("password already exists")
					}
				}
				storage.Users[i].Password = updatedUser.Password
			}

			return storage.Users[i], nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (s *UserService) DeleteUser(id int) error {
	for i, user := range storage.Users {
		if user.ID == id {
			storage.Users = append(storage.Users[:i], storage.Users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}
