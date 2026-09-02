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
	Type        *Type
}

type ExpenseModel struct {
	DB *sql.DB
}

func (m *ExpenseModel) Insert(userID int, name string, amount float64, expenseDate time.Time, typeID *int) (int, error) {
	stmt := `INSERT INTO expenses (user_id, name, amount, expense_date, created_at, type_id)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, userID, name, amount, expenseDate, typeID)
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
		SELECT e.id, e.name, e.amount, e.expense_date, t.id, t.name
		FROM expenses e
		LEFT JOIN types t ON e.type_id = t.id
		WHERE e.user_id = ?
		  AND YEAR(e.expense_date) = ?
		  AND MONTH(e.expense_date) = ?
		ORDER BY e.expense_date DESC
	`

	rows, err := m.DB.Query(stmt, userID, year, month)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var expenses []Expense

	for rows.Next() {
		var e Expense
		var typeID sql.NullInt64
		var typeName sql.NullString

		err = rows.Scan(&e.ID, &e.Name, &e.Amount, &e.ExpenseDate, &typeID, &typeName)
		if err != nil {
			return nil, err
		}

		if typeID.Valid {
			e.Type = &Type{
				ID:   int(typeID.Int64),
				Name: typeName.String,
			}
		}

		expenses = append(expenses, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (m *ExpenseModel) GetTotal(userID int, year int, month int) (float64, error) {
	stmt := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ?
		  AND YEAR(expense_date) = ?
		  AND MONTH(expense_date) = ?
	`

	var total float64

	err := m.DB.QueryRow(stmt, userID, year, month).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
