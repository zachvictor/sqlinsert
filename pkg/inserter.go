package sqlinsert

import (
	"context"
	"database/sql"
)

// Inserter models functionality to produce a valid SQL INSERT statement with bind args.
type Inserter interface {
	Tokenize(tokenType TokenType) string
	Columns() string
	Params() string
	SQL() string
	Args() []interface{}
	Insert(with InsertWith) (*sql.Stmt, error)
	InsertContext(ctx context.Context, with InsertWith) (*sql.Stmt, error)
}
