package models

import (
	"database/sql"
	"time"
)

type Expense struct {
	ID          int
	UserID      int
	Name        string
	Amount      float64
	ExpenseDate time.Time
	CreatedAt   time.Time
}

type ExpenseModel struct {
	DB *sql.DB
}

func (m *ExpenseModel) Insert(userID int, name string, amount float64, expenseDate time.Time) (int, error) {
	stmt := `INSERT INTO expenses (user_id, name, amount, expense_date, created_at)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, userID, name, amount, expenseDate)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ExpenseModel) List(userID int, year int, month int) ([]Expense, error) {
	stmt := `
		SELECT id, name, amount, expense_date
		FROM expenses
		WHERE user_id = ?
		  AND YEAR(expense_date) = ?
		  AND MONTH(expense_date) = ?
		ORDER BY expense_date DESC
	`

	rows, err := m.DB.Query(stmt, userID, year, month)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var expenses []Expense

	for rows.Next() {
		var e Expense

		err = rows.Scan(&e.ID, &e.Name, &e.Amount, &e.ExpenseDate)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}
