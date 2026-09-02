package models

import (
	"database/sql"
	"time"
)

type Income struct {
	ID         int
	UserID     int
	Name       string
	Amount     float64
	IncomeDate time.Time
	CreatedAt  time.Time
	Type       *Type
}

type IncomeModel struct {
	DB *sql.DB
}

func (m *IncomeModel) Insert(userID int, name string, amount float64, incomeDate time.Time, typeID *int) (int, error) {
	stmt := `INSERT INTO incomes(user_id, name, amount, income_date, created_at, type_id)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, userID, name, amount, incomeDate, typeID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *IncomeModel) List(userID int, year int, month int) ([]Income, error) {
	stmt := `
		SELECT i.id, i.name, i.amount, i.income_date, t.id, t.name
		FROM incomes i
		LEFT JOIN types t ON i.type_id = t.id
		WHERE i.user_id = ?
		  AND YEAR(i.income_date) = ?
		  AND MONTH(i.income_date) = ?
		ORDER BY i.income_date DESC
	`

	rows, err := m.DB.Query(stmt, userID, year, month)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var incomes []Income

	for rows.Next() {
		var i Income
		var typeID sql.NullInt64
		var typeName sql.NullString

		err = rows.Scan(&i.ID, &i.Name, &i.Amount, &i.IncomeDate, &typeID, &typeName)
		if err != nil {
			return nil, err
		}

		if typeID.Valid {
			i.Type = &Type{
				ID:   int(typeID.Int64),
				Name: typeName.String,
			}
		}

		incomes = append(incomes, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incomes, nil
}

func (m *IncomeModel) GetTotal(userID int, year int, month int) (float64, error) {
	stmt := `
		SELECT COALESCE(SUM(amount), 0)
		FROM incomes
		WHERE user_id = ?
		  AND YEAR(income_date) = ?
		  AND MONTH(income_date) = ?
	`

	var total float64

	err := m.DB.QueryRow(stmt, userID, year, month).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
