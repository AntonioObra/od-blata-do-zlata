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

// func (m *SnippetModel) Get(id int) (Snippet, error) {
// 	stmt := `SELECT id, title, content, created, expires FROM snippets
// 	WHERE expires > UTC_TIMESTAMP() AND id = ?`
//
// 	row := m.DB.QueryRow(stmt, id)
//
// 	var s Snippet
//
// 	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return Snippet{}, ErrNoRecord
// 		} else {
// 			return Snippet{}, err
// 		}
// 	}
//
// 	return s, nil
// }
//
// func (m *SnippetModel) Latest() ([]Snippet, error) {
// 	stmt := `SELECT id, title, content, created, expires FROM snippets
// 	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
//
// 	rows, err := m.DB.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer rows.Close()
//
// 	var snippets []Snippet
//
// 	for rows.Next() {
// 		var s Snippet
//
// 		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		snippets = append(snippets, s)
// 	}
//
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return snippets, nil
// }
