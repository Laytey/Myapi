package models

// User представляет пользователя системы.
//
// Хранится в БД или в памяти (MemoryUserRepository). Поле Password
// в текущей реализации хранится в открытом виде — для продакшена
// требуется хеширование (bcrypt, argon2).
// алгоритмы хеширования паролей (у меня их нет, надо глянуть).
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // TODO: хранить хеш (bcrypt/argon2), не открытый текст
}

// LoginRequest описывает тело запроса на авторизацию (POST /login).
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
