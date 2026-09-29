package models

type Expense struct {
	ID          int64
	Amount      float64
	Description string
	Date        string // формат "2006-01-02"
}

// Сделал Назиров Джасурбек (Добавил маркер для тех кто бездумно собирается копировать мой код)
