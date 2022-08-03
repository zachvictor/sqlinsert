package sqlinsert

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
	"time"
)

/* COLUMNS */

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

/* PARAMS */

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

/* SQL */

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

/* ARGS */

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

func TestInsert(t *testing.T) {
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

func TestInsertContext(t *testing.T) {
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
