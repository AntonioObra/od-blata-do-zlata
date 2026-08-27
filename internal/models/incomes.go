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
}

type IncomeModel struct {
	DB *sql.DB
}

func (m *IncomeModel) Insert(userID int, name string, amount float64, incomeDate time.Time) (int, error) {
	stmt := `INSERT INTO incomes(user_id, name, amount, income_date, created_at)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, userID, name, amount, incomeDate)
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
		SELECT id, name, amount, income_date
		FROM incomes
		WHERE user_id = ?
		  AND YEAR(income_date) = ?
		  AND MONTH(income_date) = ?
		ORDER BY income_date DESC
	`

	rows, err := m.DB.Query(stmt, userID, year, month)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var incomes []Income

	for rows.Next() {
		var i Income

		err = rows.Scan(&i.ID, &i.Name, &i.Amount, &i.IncomeDate)
		if err != nil {
			return nil, err
		}

		incomes = append(incomes, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incomes, nil
}
