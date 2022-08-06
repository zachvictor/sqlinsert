package sqlinsert

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
	"time"
)

/* Insert.Columns */

// - Single-row Insert.Columns

func TestColumnsOneRecValue(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := `(id,candy_name,form_factor,description,manufacturer,weight_grams,ts)`
	columns := ins.Columns()
	if expected != columns {
		t.Fatalf(`expected "%s", got "%s"`, expected, columns)
	}
}

func TestColumnsOneRecPointer(t *testing.T) {
	ins := Insert{tbl, recPointer}
	expected := `(id,candy_name,form_factor,description,manufacturer,weight_grams,ts)`
	columns := ins.Columns()
	if expected != columns {
		t.Fatalf(`expected "%s", got "%s"`, expected, columns)
	}
}

// - Multi-row Insert.Columns

func TestColumnsManyRecsValues(t *testing.T) {
	ins := Insert{tbl, fiveRecsValues}
	expected := `(id,candy_name,form_factor,description,manufacturer,weight_grams,ts)`
	columns := ins.Columns()
	if expected != columns {
		t.Fatalf(`expected "%s", got "%s"`, expected, columns)
	}
}

func TestColumnsManyRecsPointers(t *testing.T) {
	ins := Insert{tbl, fiveRecsPointers}
	expected := `(id,candy_name,form_factor,description,manufacturer,weight_grams,ts)`
	columns := ins.Columns()
	if expected != columns {
		t.Fatalf(`expected "%s", got "%s"`, expected, columns)
	}
}

/* Insert.Params */

// - Single-row Insert.Params

func TestParamsOneRecValue(t *testing.T) {
	UseTokenType = QuestionMarkTokenType
	ins := Insert{tbl, recValue}
	expected := `(?,?,?,?,?,?,?)`
	params := ins.Params()
	if expected != params {
		t.Fatalf(`expected "%s", got "%s"`, expected, params)
	}
}

func TestParamsOneRecPointer(t *testing.T) {
	UseTokenType = QuestionMarkTokenType
	ins := Insert{tbl, recPointer}
	expected := `(?,?,?,?,?,?,?)`
	params := ins.Params()
	if expected != params {
		t.Fatalf(`expected "%s", got "%s"`, expected, params)
	}
}

// - Multi-row Insert.Params

func TestParamsManyRecsValues(t *testing.T) {
	UseTokenType = QuestionMarkTokenType
	ins := Insert{tbl, fiveRecsValues}
	expected := `(?,?,?,?,?,?,?),(?,?,?,?,?,?,?),(?,?,?,?,?,?,?),(?,?,?,?,?,?,?),(?,?,?,?,?,?,?)`
	params := ins.Params()
	if expected != params {
		t.Fatalf(`expected "%s", got "%s"`, expected, params)
	}
}

func TestParamsManyRecsPointers(t *testing.T) {
	UseTokenType = QuestionMarkTokenType
	ins := Insert{tbl, fiveRecsPointers}
	expected := `(?,?,?,?,?,?,?),(?,?,?,?,?,?,?),(?,?,?,?,?,?,?),(?,?,?,?,?,?,?),(?,?,?,?,?,?,?)`
	params := ins.Params()
	if expected != params {
		t.Fatalf(`expected "%s", got "%s"`, expected, params)
	}
}

/* Insert.SQL */

// - Single-row Insert.SQL

func TestSQLOneRecValue(t *testing.T) {
	UseTokenType = OrdinalNumberTokenType
	ins := Insert{tbl, recValue}
	expected := `INSERT INTO candy (id,candy_name,form_factor,description,manufacturer,weight_grams,ts) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	insertSQL := ins.SQL()
	if expected != insertSQL {
		t.Fatalf(`expected "%s", got "%s"`, expected, insertSQL)
	}
}

func TestSQLOneRecPointer(t *testing.T) {
	UseTokenType = OrdinalNumberTokenType
	ins := Insert{tbl, recPointer}
	expected := `INSERT INTO candy (id,candy_name,form_factor,description,manufacturer,weight_grams,ts) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	insertSQL := ins.SQL()
	if expected != insertSQL {
		t.Fatalf(`expected "%s", got "%s"`, expected, insertSQL)
	}
}

// - Multi-row Insert.SQL

func TestSQLManyRecsValues(t *testing.T) {
	UseTokenType = OrdinalNumberTokenType
	ins := Insert{tbl, fiveRecsValues}
	expected := `INSERT INTO candy (id,candy_name,form_factor,description,manufacturer,weight_grams,ts) VALUES ($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7)`
	insertSQL := ins.SQL()
	if expected != insertSQL {
		t.Fatalf(`expected "%s", got "%s"`, expected, insertSQL)
	}
}

func TestSQLManyRecsPointers(t *testing.T) {
	UseTokenType = OrdinalNumberTokenType
	ins := Insert{tbl, fiveRecsPointers}
	expected := `INSERT INTO candy (id,candy_name,form_factor,description,manufacturer,weight_grams,ts) VALUES ($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7),($1,$2,$3,$4,$5,$6,$7)`
	insertSQL := ins.SQL()
	if expected != insertSQL {
		t.Fatalf(`expected "%s", got "%s"`, expected, insertSQL)
	}
}

/* Insert.Args */

// - Single-row Insert.Args

func TestArgsOneRecValue(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := []interface{}{
		`c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
		`Gougat`,
		`Package`,
		`tastes like gopher feed`,
		`Gouggle`,
		1.1618,
		time.Time{},
	}
	args := ins.Args()
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf(`expected "%s", got "%s"`, expected, args)
	}
}

func TestArgsOneRecPointer(t *testing.T) {
	ins := Insert{tbl, recPointer}
	expected := []interface{}{
		`c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
		`Gougat`,
		`Package`,
		`tastes like gopher feed`,
		`Gouggle`,
		1.1618,
		time.Time{},
	}
	args := ins.Args()
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf(`expected "%s", got "%s"`, expected, args)
	}
}

// - Multi-row Insert.Args

func TestArgsManyRecsValues(t *testing.T) {
	ins := Insert{tbl, fiveRecsValues}
	expected := []interface{}{
		`a`, `a`, `a`, `a`, `a`, 1.1, time.Time{},
		`b`, `b`, `b`, `b`, `b`, 2.1, time.Time{},
		`c`, `c`, `c`, `c`, `c`, 3.1, time.Time{},
		`d`, `d`, `d`, `d`, `d`, 4.1, time.Time{},
		`e`, `e`, `e`, `e`, `e`, 5.1, time.Time{},
	}
	args := ins.Args()
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf(`expected "%s", got "%s"`, expected, args)
	}
}

func TestArgsManyRecsPointers(t *testing.T) {
	ins := Insert{tbl, fiveRecsPointers}
	expected := []interface{}{
		`a`, `a`, `a`, `a`, `a`, 1.1, time.Time{},
		`b`, `b`, `b`, `b`, `b`, 2.1, time.Time{},
		`c`, `c`, `c`, `c`, `c`, 3.1, time.Time{},
		`d`, `d`, `d`, `d`, `d`, 4.1, time.Time{},
		`e`, `e`, `e`, `e`, `e`, 5.1, time.Time{},
	}
	args := ins.Args()
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf(`expected "%s", got "%s"`, expected, args)
	}
}

/* INSERT */

// - Single-row Insert.Insert, Insert.InsertContext

// TestInsertOneRecValue tests single-row insert with every token type using struct input
func TestInsertOneRecValue(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, recValue}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.Insert(db)
		if err != nil {
			t.Fatalf(`failed at Insert, could not execute SQL statement %s`, err)
		}
	}
}

// TestInsertOneRecPointer tests single-row insert with every token type using struct-pointer input
func TestInsertOneRecPointer(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, recPointer}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.Insert(db)
		if err != nil {
			t.Fatalf(`failed at Insert, could not execute SQL statement %s`, err)
		}
	}
}

// TestInsertContextOneRecValue tests single-row insert with context with every token type using struct input
func TestInsertContextOneRecValue(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, recValue}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.InsertContext(context.Background(), db)
		if err != nil {
			t.Fatalf(`failed at InsertContext, could not execute SQL statement %s`, err)
		}
	}
}

// TestInsertContextOneRecPointer tests single-row insert with context with every token type using struct-pointer input
func TestInsertContextOneRecPointer(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, recPointer}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.InsertContext(context.Background(), db)
		if err != nil {
			t.Fatalf(`failed at InsertContext, could not execute SQL statement %s`, err)
		}
	}
}

// - Multi-row Insert.Insert, Insert.InsertContext

// TestInsertManyRecsValues tests multi-row insert with every token type using slice-of-struct input
func TestInsertManyRecsValues(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, fiveRecsValues}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.Insert(db)
		if err != nil {
			t.Fatalf(`failed at Insert, could not execute SQL statement %s`, err)
		}
	}
}

// TestInsertManyRecsPointers tests multi-row insert with every token type using slice-of-struct-pointer input
func TestInsertManyRecsPointers(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, fiveRecsPointers}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.Insert(db)
		if err != nil {
			t.Fatalf(`failed at Insert, could not execute SQL statement %s`, err)
		}
	}
}

// TestInsertContextManyRecsValues tests multi-row insert with context with every token type using slice-of-struct input
func TestInsertContextManyRecsValues(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, fiveRecsValues}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.InsertContext(context.Background(), db)
		if err != nil {
			t.Fatalf(`failed at InsertContext, could not execute SQL statement %s`, err)
		}
	}
}

// TestInsertContextManyRecsPointers tests multi-row insert with context with every token type using slice-of-struct-pointer input
func TestInsertContextManyRecsPointers(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := Insert{tbl, fiveRecsPointers}
		s := regexp.QuoteMeta(ins.SQL())
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf(`failed to construct SQL mock %s`, err)
		}
		mock.ExpectPrepare(s)
		mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
		_, err = ins.InsertContext(context.Background(), db)
		if err != nil {
			t.Fatalf(`failed at InsertContext, could not execute SQL statement %s`, err)
		}
	}
}
