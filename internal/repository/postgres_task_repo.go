package repository

import (
	"Myapi/internal/db"
	"Myapi/internal/models"
	"context"
	"errors"
)

type PostgresTaskRepository struct {
	storage *db.Storage
}

func NewPostgresTaskRepository(storage *db.Storage) *PostgresTaskRepository {
	return &PostgresTaskRepository{storage: storage}
}

func (r *PostgresTaskRepository) Save(task *models.Task) error {
	query := `INSERT INTO tasks (title, description, status, user_uid) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.storage.Pool.QueryRow(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UserUID,
	).Scan(&task.ID)
	return err
}

func (r *PostgresTaskRepository) GetByID(id int) (models.Task, error) {
	query := `SELECT id, title, description, status, user_uid FROM tasks WHERE id = $1 AND deleted = false`

	var task models.Task
	err := r.storage.Pool.QueryRow(context.Background(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserUID,
	)
	if err != nil {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

func (r *PostgresTaskRepository) GetByUserUID(uid string) ([]models.Task, error) {
	query := `SELECT id, title, description, status, user_uid FROM tasks WHERE user_uid = $1 AND deleted = false`

	rows, err := r.storage.Pool.Query(context.Background(), query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserUID); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) GetAll() ([]models.Task, error) {
	query := `SELECT id, title, description, status, user_uid FROM tasks WHERE deleted = false`

	rows, err := r.storage.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserUID); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) Update(task models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, user_uid = $4 WHERE id = $5`
	_, err := r.storage.Pool.Exec(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UserUID,
		task.ID,
	)
	return err
}

func (r *PostgresTaskRepository) Delete(id int) error {
	_, err := r.storage.Pool.Exec(
		context.Background(),
		`UPDATE tasks SET deleted = true WHERE id = $1`,
		id,
	)
	return err
}

// Scan:
// Получает данные из результата запроса.
// Преобразует их в нужный тип.
// Записывает в переменную, на которую указывает переданный указатель.

func (r *PostgresTaskRepository) HardDelete() error {
	tx, err := r.storage.Pool.Begin(context.Background())
	// Begin берёт соединение из пула и помечает его как «в транзакции»
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	// Rollback(ctx) закрывает соединение и отменяет транзакцию

	_, err = tx.Exec(context.Background(), `DELETE FROM tasks WHERE deleted = true`)
	// Exec выполняет SQL-запрос и возвращает количество измененных строк
	// Exec используется для запросов, которые не возвращают данные (INSERT UPDATE DELETE)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
	// Commit(ctx) фиксирует все изменения, сделанные внутри транзакции

}
