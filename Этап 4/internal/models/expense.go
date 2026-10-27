package models

type Category struct {
	ID   int64
	Name string
}

type Expense struct {
	ID           int64
	Amount       float64
	Description  string
	Date         string // формат "2006-01-02"
	CategoryID   int64
	CategoryName string // для вывода в списке через JOIN
}
