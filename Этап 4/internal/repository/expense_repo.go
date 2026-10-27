package repository

import (
	"database/sql"

	"NazirovJasur/my-http-server.git/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// --- РАБОТА С КАТЕГОРИЯМИ ---

func (r *ExpenseRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

// --- РАБОТА С ТРАТАМИ ---

func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		`INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)`,
		e.Amount, e.Description, e.Date, e.CategoryID,
	)
	return err
}

func (r *ExpenseRepository) GetAll(categoryID int64, dateFrom, dateTo string) ([]models.Expense, error) {
	// Динамическое построение SQL-условий
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name 
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE 1=1
	`
	var args []interface{}

	if categoryID > 0 {
		query += " AND e.category_id = ?"
		args = append(args, categoryID)
	}
	if dateFrom != "" {
		query += " AND e.date >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND e.date <= ?"
		args = append(args, dateTo)
	}

	query += " ORDER BY e.date DESC, e.id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *ExpenseRepository) GetByID(id int64) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		`SELECT id, amount, description, date, category_id FROM expenses WHERE id = ?`, id,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID)
	return e, err
}

func (r *ExpenseRepository) Update(e models.Expense) error {
	res, err := r.db.Exec(
		`UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?`,
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ExpenseRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM expenses WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
