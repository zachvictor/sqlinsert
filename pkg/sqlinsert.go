package sqlinsert

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// UseStructTag specifies the struct tag key for the column name. Default is `col`.
var UseStructTag = `col`

// TokenType represents a type of token in a SQL INSERT statement, whether column or value expression.
type TokenType int

const (

	/* COLUMN TokenType */

	// ColumnNameTokenType uses the column name from the struct tag specified by UseStructTag.
	// INSERT INTO tbl (foo, bar, ... baz)
	ColumnNameTokenType TokenType = 0

	/* VALUE TokenType */

	// QuestionMarkTokenType uses question marks as value-tokens.
	// VALUES (?, ?, ... ?) -- MySQL, SingleStore
	QuestionMarkTokenType TokenType = 1

	// AtColumnNameTokenType uses @ followed by the column name from the struct tag specified by UseStructTag.
	// VALUES (@foo, @bar, ... @baz) -- MySQL, SingleStore
	AtColumnNameTokenType TokenType = 2

	// OrdinalNumberTokenType uses % plus the value of an ordered sequence of integers starting at 1.
	// %1, %2, ... %n -- Postgres
	OrdinalNumberTokenType TokenType = 3

	// ColonTokenType uses : followed by the column name from the struct tag specified by UseStructTag.
	// :foo, :bar, ... :baz -- Oracle
	ColonTokenType TokenType = 4
)

// UseTokenType specifies the token type to use for values. Default is the question mark (`?`).
var UseTokenType = QuestionMarkTokenType

// InsertWith models functionality needed to execute a SQL INSERT statement with database/sql via sql.DB or sql.Tx.
// Note: sql.Conn is also supported, however, for PrepareContext and ExecContext only.
type InsertWith interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

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

// Insert models data used to produce a valid SQL INSERT statement with bind args.
// Table is the table name. Record is the struct with column-name tagged fields and the data to be inserted.
// Private fields recordType and recordValue are used internally with reflection and interfaces to field values.
type Insert struct {
	Table       string
	Record      interface{}
	recordType  reflect.Type
	recordValue reflect.Value
}

// NewInsert builds a new Insert using a given table-name (string) and record data (struct).
// It is recommended that new NewInsert be used to build every new INSERT; however, if only Insert.Tokenize is
// needed, a "manually" build Insert will support tokenization as long as Insert.Table and Insert.Record are valid.
func NewInsert(Table string, Record interface{}) *Insert {
	return &Insert{
		Table:       Table,
		Record:      Record,
		recordType:  reflect.TypeOf(Record),
		recordValue: reflect.ValueOf(Record),
	}
}

// Tokenize translates struct fields into the tokens of SQL column or value expressions.
func (ins *Insert) Tokenize(tokenType TokenType) string {
	var b strings.Builder
	for i := 0; i < ins.recordType.NumField(); i++ {
		switch tokenType {
		case ColumnNameTokenType:
			b.WriteString(ins.recordType.Field(i).Tag.Get(UseStructTag))
		case QuestionMarkTokenType:
			_, _ = fmt.Fprint(&b, `?`)
		case AtColumnNameTokenType:
			_, _ = fmt.Fprintf(&b, `@%s`, ins.recordType.Field(i).Tag.Get(UseStructTag))
		case OrdinalNumberTokenType:
			_, _ = fmt.Fprintf(&b, `$%d`, i+1)
		case ColonTokenType:
			_, _ = fmt.Fprintf(&b, `:%s`, ins.recordType.Field(i).Tag.Get(UseStructTag))
		}
		if i < ins.recordType.NumField()-1 {
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
func (ins *Insert) Params() string {
	return ins.Tokenize(UseTokenType)
}

// SQL returns the full parameterized SQL INSERT statement.
func (ins *Insert) SQL() string {
	var insertSQL strings.Builder
	_, _ = fmt.Fprintf(&insertSQL, `INSERT INTO %s (%s) VALUES (%s)`,
		ins.Table, ins.Columns(), ins.Params())
	return insertSQL.String()
}

// Args returns the arguments to be bound in Insert() or the variadic Exec/ExecContext functions in database/sql.
func (ins *Insert) Args() []interface{} {
	args := make([]interface{}, ins.recordType.NumField())
	for i := 0; i < ins.recordType.NumField(); i++ {
		args[i] = ins.recordValue.Field(i).Interface()
	}
	return args
}

// Insert prepares and executes a SQL INSERT statement on a *sql.DB, *sql.Tx,
// or other Inserter-compatible interface to Prepare and Exec.
func (ins *Insert) Insert(with InsertWith) (*sql.Stmt, error) {
	stmt, err := with.Prepare(ins.SQL())
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(ins.Args()...)
	return stmt, err
}

// InsertContext prepares and executes a SQL INSERT statement on a *sql.DB, *sql.Tx, *sql.Conn,
// or other Inserter-compatible interface to PrepareContext and ExecContext.
func (ins *Insert) InsertContext(ctx context.Context, with InsertWith) (*sql.Stmt, error) {
	stmt, err := with.Prepare(ins.SQL())
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.ExecContext(ctx, ins.Args()...)
	return stmt, err
}
