package models

import "database/sql"

type Type struct {
	ID     int
	UserID int
	Name   string
}

type TypeModel struct {
	DB *sql.DB
}

func (m *TypeModel) Insert(userID int, name string) (int, error) {
	stmt := `INSERT INTO types(user_id, name) VALUES(?, ?)`

	result, err := m.DB.Exec(stmt, userID, name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TypeModel) List(userID int) ([]Type, error) {
	stmt := `SELECT id, name FROM types WHERE user_id = ? ORDER BY id DESC`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var types []Type

	for rows.Next() {
		var t Type

		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}

		types = append(types, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}
