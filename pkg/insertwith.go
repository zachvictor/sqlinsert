package sqlinsert

import (
	"context"
	"database/sql"
)

// InsertWith models functionality needed to execute a SQL INSERT statement with database/sql via sql.DB or sql.Tx.
// Note: sql.Conn is also supported, however, for PrepareContext and ExecContext only.
type InsertWith interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
