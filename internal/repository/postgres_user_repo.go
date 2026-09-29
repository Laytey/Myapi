package repository

import (
	"context"
	"errors"

	"Myapi/internal/db"
	"Myapi/internal/models"
)

// PostgresUserRepository — реализация UserRepository на PostgreSQL
// через пул соединений pgxpool.
//
// Потокобезопасна: каждое соединение берётся из пула и не используется
// параллельно двумя горутинами. Данные сохраняются между перезапусками.
type PostgresUserRepository struct {
	storage *db.Storage
}

// NewPostgresUserRepository создаёт репозиторий, использующий
// переданный пул соединений с PostgreSQL.
func NewPostgresUserRepository(storage *db.Storage) *PostgresUserRepository {
	return &PostgresUserRepository{storage: storage}
}

// Save реализует UserRepository.Save.
// Вставляет нового пользователя в таблицу users. ID генерируется БД
// через SERIAL и возвращается через RETURNING id.
func (r *PostgresUserRepository) Save(user *models.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	// $1, $2, $3 — это плейсхолдеры (места для подстановки). Они защищают от SQL-инъекций
	// TODO: разобраться с SQL-инъекциями
	// QueryRow выполняет запрос и возвращает одну строку результата (в этом случае id)
	err := r.storage.Pool.QueryRow(context.Background(), query, user.Name, user.Email, user.Password).Scan(&user.ID)
	//context.Background() создаёт пустой контекст
	// Scan записывает результат запроса (id) в поле user.ID, меняет значение
	return err
}

// GetByID реализует UserRepository.GetByID.
// Возвращает ошибку "user not found", если пользователь не существует.
func (r *PostgresUserRepository) GetByID(id int) (models.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`
	var user models.User // выделение места под данные
	// создаём новую переменную, которую потом заполним данными из базы
	err := r.storage.Pool.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// GetByEmail реализует UserRepository.GetByEmail.
// Используется при логине и проверке уникальности email при регистрации.
// Возвращает ошибку "user not found", если пользователь с таким email
// не зарегистрирован.
func (r *PostgresUserRepository) GetByEmail(email string) (models.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	var user models.User
	err := r.storage.Pool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// GetAll реализует UserRepository.GetAll.
// Возвращает всех пользователей без фильтрации.
func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	rows, err := r.storage.Pool.Query(context.Background(), `SELECT id, name, email, password FROM users`)
	// Query ожидает несколько строк, в отличие от QueryRow, который ждёт одну строку

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		// Next переходит к следующей строке результата запроса и читает до тех пор, пока они есть
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Update реализует UserRepository.Update.
// Обновляет все поля пользователя, кроме ID.
func (r *PostgresUserRepository) Update(user models.User) error {
	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`
	_, err := r.storage.Pool.Exec(context.Background(), query, user.Name, user.Email, user.Password, user.ID)
	// Exec выполняет SQL-запрос и возвращает количество измененных строк
	// Exec используется для запросов, которые не возвращают данные (INSERT UPDATE DELETE)
	return err
}

// Delete реализует UserRepository.Delete.
// В отличие от PostgresTaskRepository, выполняет физическое удаление:
// строка полностью удаляется из таблицы users.
func (r *PostgresUserRepository) Delete(id int) error {
	_, err := r.storage.Pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	return err
}
