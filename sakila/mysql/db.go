package mysql

import (
	"context"
	"database/sql"
	"strings"
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

func argString(count int) string {
	argStrings := make([]string, count)
	for i := 0; i < count; i++ {
		argStrings[i] = "?"
	}

	return strings.Join(argStrings, ", ")
}
