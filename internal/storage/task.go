package storage

import (
	"Myapi/internal/models"
)

// временное хранилище задач в памяти
var Tasks = []models.Task{
	{ID: 1, Title: "Изучить тему", Description: "Выполнить практику", Status: "Новая"},
	{ID: 2, Title: "Дополнить код", Description: "Добавить методы", Status: "В процессе"},
}
