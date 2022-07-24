package sqlinsert

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// StructTag specifies the struct tag key for the column name. Default is `col`.
var StructTag = `col`

// TokenType represents a token in a SQL statement, whether column or value expression.
type TokenType int

const (

	/* COLUMN TokenType */

	// ColumnNameTokenType uses the column name from the struct tag specified by StructTag.
	// INSERT INTO tbl (foo, bar, ... baz)
	ColumnNameTokenType TokenType = 0

	/* VALUE TokenType */

	// QuestionMarkTokenType uses question marks as value-tokens.
	// VALUES (?, ?, ... ?) -- MySQL, SingleStore
	QuestionMarkTokenType = 1

	// AtColumnNameTokenType uses @ followed by the column name from the struct tag specified by StructTag.
	// VALUES (@foo, @bar, ... @baz) -- MySQL, SingleStore
	AtColumnNameTokenType = 2

	// OrdinalNumberTokenType uses % plus the value of an ordered sequence of integers starting at 1.
	// %1, %2, ... %n -- Postgres
	OrdinalNumberTokenType = 3

	// ColonTokenType uses : followed by the column name from the struct tag specified by StructTag.
	// :foo, :bar, ... :baz -- Oracle
	ColonTokenType = 4
)

// InsertWith models functionality needed to execute a SQL INSERT statement with database/sql via Conn, DB, or Tx.
type InsertWith interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// Inserter models functionality to produce a valid SQL INSERT statement with bind args.
type Inserter interface {
	Tokenize(tokenType TokenType) string
	Columns() string
	Params(tokenType TokenType) string
	SQL(tokenType TokenType) string
	Args() []interface{}
	Insert(tokenType TokenType, with InsertWith) (*sql.Stmt, error)
	InsertContext(tokenType TokenType, with InsertWith) (*sql.Stmt, error)
}

// Insert models data used to produce a valid SQL INSERT statement with bind args.
type Insert struct {
	Table    string
	Record   interface{}
	RowType  reflect.Type
	RowValue reflect.Value
}

// NewInsert builds a new Insert using a given table-name (string) and row-data (struct).
func NewInsert(Table string, Record interface{}) *Insert {
	return &Insert{
		Table:    Table,
		Record:   Record,
		RowType:  reflect.TypeOf(Record),
		RowValue: reflect.ValueOf(Record),
	}
}

// Tokenize translates struct fields into the tokens of SQL column or value expressions.
func (ins *Insert) Tokenize(tokenType TokenType) string {
	var b strings.Builder
	for i := 0; i < ins.RowType.NumField(); i++ {
		switch tokenType {
		case ColumnNameTokenType:
			b.WriteString(ins.RowType.Field(i).Tag.Get(StructTag))
		case QuestionMarkTokenType:
			_, _ = fmt.Fprint(&b, `?`)
		case AtColumnNameTokenType:
			_, _ = fmt.Fprintf(&b, `@%s`, ins.RowType.Field(i).Tag.Get(StructTag))
		case OrdinalNumberTokenType:
			_, _ = fmt.Fprintf(&b, `$%d`, i+1)
		case ColonTokenType:
			_, _ = fmt.Fprintf(&b, `:%s`, ins.RowType.Field(i).Tag.Get(StructTag))
		}
		if i < ins.RowType.NumField()-1 {
			b.WriteString(`, `)
		}
	}
	return b.String()
}

// Columns returns the comma-separated list of column names-as-tokens for the SQL INSERT statement.
func (ins *Insert) Columns() string {
	return ins.Tokenize(ColumnNameTokenType)
}

// Params returns the comma-separated list of bind param tokens for the SQL INSERT statement.
func (ins *Insert) Params(tokenType TokenType) string {
	return ins.Tokenize(tokenType)
}

// SQL returns the full parameterized SQL INSERT statement.
func (ins *Insert) SQL(tokenType TokenType) string {
	var insertSQL strings.Builder
	_, _ = fmt.Fprintf(&insertSQL, `INSERT INTO %s (%s) VALUES (%s)`,
		ins.Table, ins.Columns(), ins.Params(tokenType))
	return insertSQL.String()
}

// Args returns the arguments to be bound in Insert() or the variadic Exec/ExecContext functions in database/sql.
func (ins *Insert) Args() []interface{} {
	args := make([]interface{}, ins.RowType.NumField())
	for i := 0; i < ins.RowType.NumField(); i++ {
		args[i] = ins.RowValue.Field(i).Interface()
	}
	return args
}

// Insert prepares and executes a SQL INSERT statement on a *sql.Conn, *sql.DB, *sql.Tx, or equivalent interface supporting Prepare(string).
func (ins *Insert) Insert(tokenType TokenType, with InsertWith) (*sql.Stmt, error) {
	stmt, err := with.Prepare(ins.SQL(tokenType))
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(ins.Args()...)
	return stmt, err
}

// InsertContext prepares and executes a SQL INSERT statement on a *sql.Conn, *sql.DB, *sql.Tx, or equivalent interface supporting PrepareContext(context.Context, string).
func (ins *Insert) InsertContext(tokenType TokenType, with InsertWith) (*sql.Stmt, error) {
	stmt, err := with.Prepare(ins.SQL(tokenType))
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(ins.Args()...)
	return stmt, err
}
