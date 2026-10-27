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

func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		`INSERT INTO expenses (user_id, amount, description, date, category_id) VALUES (?, ?, ?, ?, ?)`,
		e.UserID, e.Amount, e.Description, e.Date, e.CategoryID,
	)
	return err
}

// GetAll теперь фильтрует по user_id
func (r *ExpenseRepository) GetAll(userID int64, categoryID int64, dateFrom, dateTo string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.user_id, e.amount, e.description, e.date, e.category_id, c.name 
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?
	`
	args := []interface{}{userID}

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
		if err := rows.Scan(&e.ID, &e.UserID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *ExpenseRepository) GetByID(userID, id int64) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		`SELECT id, user_id, amount, description, date, category_id FROM expenses WHERE id = ? AND user_id = ?`, id, userID,
	).Scan(&e.ID, &e.UserID, &e.Amount, &e.Description, &e.Date, &e.CategoryID)
	return e, err
}

func (r *ExpenseRepository) Update(e models.Expense) error {
	res, err := r.db.Exec(
		`UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?`,
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID, e.UserID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ExpenseRepository) Delete(userID, id int64) error {
	res, err := r.db.Exec(`DELETE FROM expenses WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
