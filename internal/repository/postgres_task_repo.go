package repository

import (
	"context"
	"errors"

	"Myapi/internal/db"
	"Myapi/internal/models"
)

// PostgresTaskRepository — реализация TaskRepository на PostgreSQL
// через пул соединений pgxpool.
//
// Потокобезопасна: каждое соединение берётся из пула и не используется
// параллельно двумя горутинами. Данные сохраняются между перезапусками.
type PostgresTaskRepository struct {
	storage *db.Storage
}

// NewPostgresTaskRepository создаёт репозиторий, использующий
// переданный пул соединений с PostgreSQL.
func NewPostgresTaskRepository(storage *db.Storage) *PostgresTaskRepository {
	return &PostgresTaskRepository{storage: storage}
}

// Save реализует TaskRepository.Save.
// Вставляет новую строку в таблицу tasks. ID генерируется БД
// через SERIAL и возвращается через RETURNING id.
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

// GetByID реализует TaskRepository.GetByID.
// Возвращает ошибку "task not found", если задача не найдена
// или помечена удалённой (deleted = true).
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

// GetByUserUID реализует TaskRepository.GetByUserUID.
// Возвращает только активные задачи (deleted = false).
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

// GetAll реализует TaskRepository.GetAll.
// Возвращает только активные задачи (deleted = false).
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

// Update реализует TaskRepository.Update.
// Обновляет все поля задачи, кроме ID и deleted.
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

// Delete реализует TaskRepository.Delete.
// Выполняет soft-delete: UPDATE tasks SET deleted = true.
// Задача остаётся в таблице, но не возвращается в GetByID/GetAll.
func (r *PostgresTaskRepository) Delete(id int) error {
	_, err := r.storage.Pool.Exec(
		context.Background(),
		`UPDATE tasks SET deleted = true WHERE id = $1`,
		id,
	)
	return err
}

// HardDelete реализует TaskRepository.HardDelete.
// Физически удаляет все задачи с deleted = true в одной транзакции.
// При ошибке — откат через Rollback (defer).
func (r *PostgresTaskRepository) HardDelete() error {
	tx, err := r.storage.Pool.Begin(context.Background())
	// Begin берёт соединение из пула и помечает его как «в транзакции»
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()
	// Rollback отменяет транзакцию и возвращает соединение в пул

	_, err = tx.Exec(context.Background(), `DELETE FROM tasks WHERE deleted = true`)
	// Exec выполняет SQL-запрос и возвращает количество измененных строк
	// Exec используется для запросов, которые не возвращают данные (INSERT UPDATE DELETE)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
	// Commit(ctx) фиксирует все изменения, сделанные внутри транзакции

}
