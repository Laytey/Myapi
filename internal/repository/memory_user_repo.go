package repository

import (
	"Myapi/internal/models"
	"errors"
	"fmt"
)

type MemoryUserRepository struct {
	users  []models.User
	nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{ //
		users:  []models.User{}, // пустой слайс
		nextID: 1,
	}
}

func (r *MemoryUserRepository) Save(user *models.User) error {
	user.ID = r.nextID
	r.nextID++
	r.users = append(r.users, *user)
	fmt.Println("Save: user.ID =", user.ID)
	return nil // возвращает error.nil
}

func (r *MemoryUserRepository) GetByID(id int) (models.User, error) {
	for _, u := range r.users { // u -каждый пользователь в цикле
		if u.ID == id {
			return u, nil
		}
	} // models.User{} - пустой пользователь
	return models.User{}, errors.New("user not found")
}

func (r *MemoryUserRepository) GetByEmail(email string) (models.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (r *MemoryUserRepository) GetAll() ([]models.User, error) {
	return r.users, nil
}

func (r *MemoryUserRepository) Update(user models.User) error {
	for i, u := range r.users { // user.ID - ID пользователя, кот мы обновляем
		if u.ID == user.ID { //u.ID - ID пользователя из хранилища
			r.users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *MemoryUserRepository) Delete(id int) error {
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}
