// Package models содержит структуры данных, используемые в API Myapi.
package models

// Task представляет задачу пользователя.
//
// Каждая задача привязана к конкретному пользователю через UserUID.
// Поле Deleted используется для soft-delete: задача не удаляется
// физически, а помечается как удалённая.
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserUID     string `json:"user_uid"`
	Deleted     bool   `json:"deleted"`
}
