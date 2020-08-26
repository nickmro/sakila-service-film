package mysql

import "database/sql"

// Open opens a MySQL database connection.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}
