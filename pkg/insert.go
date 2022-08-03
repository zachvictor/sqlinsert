package sqlinsert

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

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
		args []interface{}
		t    reflect.Type
		row  reflect.Value
	)
	v := reflect.ValueOf(ins.Data)
	if v.Kind() == reflect.Slice {
		argIndex := -1
		if v.Index(0).Kind() == reflect.Pointer {
			t = v.Index(0).Elem().Type()
		} else {
			t = v.Index(0).Type()
		}
		args = make([]interface{}, v.Len()*t.NumField())
		for rowIndex := 0; rowIndex < v.Len(); rowIndex++ {
			for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
				argIndex += 1
				if v.Index(0).Kind() == reflect.Pointer {
					row = v.Index(rowIndex).Elem()
				} else {
					row = v.Index(rowIndex)
				}
				args[argIndex] = row.Field(fieldIndex).Interface()
			}
		}
		return args
	} else {
		if v.Kind() == reflect.Pointer {
			t = v.Elem().Type()
			row = v.Elem()
		} else {
			t = v.Type()
			row = v
		}
		args = make([]interface{}, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			args[i] = row.Field(i).Interface()
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
