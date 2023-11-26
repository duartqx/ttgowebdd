package sqlite

import "github.com/jmoiron/sqlx"

func NewDbConnection(connStr string) (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", connStr)
}
