package service

import (
	"Myapi/internal/models"
	"Myapi/internal/repository"
	"testing"
)

func TestCreateTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	svc := NewTaskServiceForTest(repo)

	_, err := svc.CreateTask(models.Task{Title: ""})
	if err == nil {
		t.Errorf("ожидали ошибку при пустом Title, получили nil")
	}

	task, err := svc.CreateTask(models.Task{Title: "Test task"})
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if task.Status != "Новая" {
		t.Errorf("ожидали статус 'Новая', получили %q", task.Status)
	}
	if task.ID != 1 {
		t.Errorf("ожидали ID=1, получили %d", task.ID)
	}

	// Кейс 3: Status задан явно - не перезаписывается
	task2, err := svc.CreateTask(models.Task{Title: "Test task 2", Status: "In Progress"})
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if task2.Status != "In Progress" {
		t.Errorf("ожидали Status='In Progress', получили %q", task2.Status)
	}

}

func TestDeleteTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	svc := NewTaskServiceForTest(repo)

	// Кейс 1: удаление несуществующей задачи - ошибка, канал пуст

	err := svc.DeleteTask(999)
	if err == nil {
		t.Errorf("ожидали ошибку 'task not found', получили nil")
	}
	if len(svc.cleanupCh) != 0 {
		t.Errorf("канал должен быть пуст, получили %d сигналов", len(svc.cleanupCh))

	}

	// Кейс 2: удаление существующей задачи - успех, задача скрыта, 1 сигнал
	task, err := svc.CreateTask(models.Task{Title: "Test task"})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	err = svc.DeleteTask(task.ID)
	if err != nil {
		t.Errorf("не ожидали ошибку при удалении: %v", err)
	}
	// Проверка: задача больше не видна через GetByID
	_, err = repo.GetByID(task.ID)
	if err == nil {
		t.Errorf("задача должна быть скрыта после DeleteTask, но GetByID её нашёл")
	}

	// Проверка: в канале появился сигнал
	if len(svc.cleanupCh) != 1 {
		t.Errorf("ожидали 1 сигнал в канале, получили %d", len(svc.cleanupCh))
	}
}

func TestGetTaskByID(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	svc := NewTaskServiceForTest(repo)

	// Кейс 1: несуществующий ID - ошибка
	_, err := svc.GetTaskByID(999)
	if err == nil {
		t.Errorf("ожидали ошибку 'task not found', получили nil")
	}

	// Кейс 2: существующий ID - задача найдена
	created, err := svc.CreateTask(models.Task{Title: "Test"})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	found, err := svc.GetTaskByID(created.ID)
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if found.Title != "Test" {
		t.Errorf("ожидали Title='Test', получили %q", found.Title)
	}
}

func TestUpdateTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	svc := NewTaskServiceForTest(repo)

	// Кейс 1: несуществующая задача - ошибка
	err := svc.UpdateTask(models.Task{ID: 999, Title: "X"})
	if err == nil {
		t.Errorf("ожидали ошибку 'task not found', получили nil")
	}

	// Кейс 2: существующая задача - успех
	created, err := svc.CreateTask(models.Task{Title: "Original"})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	updated := models.Task{
		ID:     created.ID,
		Title:  "Updated",
		Status: "Done",
	}
	err = svc.UpdateTask(updated)
	if err != nil {
		t.Errorf("не ожидали ошибку при обновлении: %v", err)
	}

	// проверяем, что поля реально изменились
	found, err := svc.GetTaskByID(created.ID)
	if err != nil {
		t.Fatalf("не удалось получить задачу после обновления: %v", err)
	}
	if found.Title != "Updated" {
		t.Errorf("ожидали Title='Updated', получили %q", found.Title)
	}
	if found.Status != "Done" {
		t.Errorf("ожидали Status='Done', получили %q", found.Status)
	}
}

func TestGetTasksByUser(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	svc := NewTaskServiceForTest(repo)

	// создаём задачи для двух пользователей
	if _, err := svc.CreateTask(models.Task{Title: "User1 Task", UserUID: "user1"}); err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}
	if _, err := svc.CreateTask(models.Task{Title: "User2 Task", UserUID: "user2"}); err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}
	if _, err := svc.CreateTask(models.Task{Title: "User1 Task 2", UserUID: "user1"}); err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	// запрашиваем задачи user1
	tasks, err := svc.GetTasksByUser("user1")
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("ожидали 2 задачи для user1, получили %d", len(tasks))
	}
	// проверяем, что вернулись именно задачи user1
	for _, task := range tasks {
		if task.UserUID != "user1" {
			t.Errorf("ожидали только задачи user1, получили задачу с UserUID=%q", task.UserUID)
		}
	}
}

func TestGetAllTasks(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	svc := NewTaskServiceForTest(repo)

	// пустой список
	tasks, err := svc.GetAllTasks()
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("ожидали 0 задач, получили %d", len(tasks))
	}

	// создаём 3 задачи
	if _, err := svc.CreateTask(models.Task{Title: "Task 1", UserUID: "user1"}); err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}
	if _, err := svc.CreateTask(models.Task{Title: "Task 2", UserUID: "user2"}); err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}
	if _, err := svc.CreateTask(models.Task{Title: "Task 3", UserUID: "user1"}); err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	tasks, err = svc.GetAllTasks()
	if err != nil {
		t.Errorf("не ожидали ошибку: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("ожидали 3 задачи, получили %d", len(tasks))
	}
}
