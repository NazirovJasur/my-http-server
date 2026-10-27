package models

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

type Category struct {
	ID   int64
	Name string
}

type Expense struct {
	ID           int64
	UserID       int64 // Привязка к пользователю
	Amount       float64
	Description  string
	Date         string
	CategoryID   int64
	CategoryName string
}
