package mysql

import (
	"context"
	"database/sql"
	"time"
)

// DB is a SQL DB connection.
type DB struct {
	*sql.DB
}

const pingTimeoutDuration = time.Second * 5

// Status satisfies the health checker interface.
func (db *DB) Status() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutDuration)
	defer cancel()

	return nil, db.DB.PingContext(ctx)
}
