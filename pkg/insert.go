package sqlinsert

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

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
// Table is the table name. Data is either a struct with column-name tagged fields and the data to be inserted or
// a slice struct (struct ptr works too). Private recordType and recordValue fields are used with reflection to get
// struct tags for Insert.Columns, Insert.Params, and Insert.SQL and to retrieve values for Insert.Args.
type Insert struct {
	Table string
	Data  interface{}
}

// Columns returns the comma-separated list of column names-as-tokens for the SQL INSERT statement.
// Multi Row Insert: Insert.Data is a slice; first item in slice is
func (ins *Insert) Columns() string {
	v := reflect.ValueOf(ins.Data)
	if v.Kind() == reflect.Slice {
		if v.Index(0).Kind() == reflect.Pointer {
			return Tokenize(v.Index(0).Elem().Type(), ColumnNameTokenType)
		} else {
			return Tokenize(v.Index(0).Type(), ColumnNameTokenType)
		}
	} else if v.Kind() == reflect.Pointer {
		return Tokenize(v.Elem().Type(), ColumnNameTokenType)
	} else {
		return Tokenize(v.Type(), ColumnNameTokenType)
	}
}

// Params returns the comma-separated list of bind param tokens for the SQL INSERT statement.
func (ins *Insert) Params() string {
	v := reflect.ValueOf(ins.Data)
	if v.Kind() == reflect.Slice {
		var (
			b        strings.Builder
			paramRow string
		)
		if v.Index(0).Kind() == reflect.Pointer {
			paramRow = Tokenize(v.Index(0).Elem().Type(), UseTokenType)
		} else {
			paramRow = Tokenize(v.Index(0).Type(), UseTokenType)
		}
		b.WriteString(paramRow)
		for i := 1; i < v.Len(); i++ {
			b.WriteString(`,`)
			b.WriteString(paramRow)
		}
		return b.String()
	} else if v.Kind() == reflect.Pointer {
		return Tokenize(v.Elem().Type(), UseTokenType)
	} else {
		return Tokenize(v.Type(), UseTokenType)
	}
}

// SQL returns the full parameterized SQL INSERT statement.
func (ins *Insert) SQL() string {
	var insertSQL strings.Builder
	_, _ = fmt.Fprintf(&insertSQL, `INSERT INTO %s %s VALUES %s`,
		ins.Table, ins.Columns(), ins.Params())
	return insertSQL.String()
}

// Args returns the arguments to be bound in Insert() or the variadic Exec/ExecContext functions in database/sql.
func (ins *Insert) Args() []interface{} {
	var (
		data    reflect.Value
		rec     reflect.Value
		recType reflect.Type
		args    []interface{}
	)
	data = reflect.ValueOf(ins.Data)
	if data.Kind() == reflect.Slice { // Multi row INSERT: Insert.Data is a slice-of-struct-pointer or slice-of-struct
		argIndex := -1
		if data.Index(0).Kind() == reflect.Pointer { // First slice element is struct pointers
			recType = data.Index(0).Elem().Type()
		} else { // First slice element is struct
			recType = data.Index(0).Type()
		}
		numRecs := data.Len()
		numFieldsPerRec := recType.NumField()
		numBindArgs := numRecs * numFieldsPerRec
		args = make([]interface{}, numBindArgs)
		for rowIndex := 0; rowIndex < data.Len(); rowIndex++ {
			if data.Index(0).Kind() == reflect.Pointer {
				rec = data.Index(rowIndex).Elem() // Cur slice elem is struct pointer, get arg val from ref-element
			} else {
				rec = data.Index(rowIndex) // Cur slice elem is struct, can get arg val directly
			}
			for fieldIndex := 0; fieldIndex < numFieldsPerRec; fieldIndex++ {
				argIndex += 1
				args[argIndex] = rec.Field(fieldIndex).Interface()
			}
		}
		return args
	} else { // Single-row INSERT: Insert.Data must be a struct pointer or struct (otherwise reflect will panic)
		if data.Kind() == reflect.Pointer { // Row information via struct pointer
			recType = data.Elem().Type()
			rec = data.Elem()
		} else { // Row information via struct
			recType = data.Type()
			rec = data
		}
		args = make([]interface{}, recType.NumField())
		for i := 0; i < recType.NumField(); i++ {
			args[i] = rec.Field(i).Interface()
		}
		return args
	}
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
