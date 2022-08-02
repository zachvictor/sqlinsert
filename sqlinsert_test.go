package sqlinsert

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
	"time"
)

type candyInsert struct {
	Id          string    `col:"id"`
	Name        string    `col:"candy_name"`
	FormFactor  string    `col:"form_factor"`
	Description string    `col:"description"`
	Mfr         string    `col:"manufacturer"`
	Weight      float64   `col:"weight_grams"`
	Timestamp   time.Time `col:"ts"`
}

var tbl = `candy`

var rec = candyInsert{
	Id:          `c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
	Name:        `Gougat`,
	FormFactor:  `Package`,
	Description: `tastes like gopher feed`,
	Mfr:         `Gouggle`,
	Weight:      1.16180,
	Timestamp:   time.Time{},
}

var valuesTokenTypes = []TokenType{
	QuestionMarkTokenType,
	AtColumnNameTokenType,
	OrdinalNumberTokenType,
	ColonTokenType,
}

type failingMock struct {
}

func (*failingMock) Prepare(query string) (*sql.Stmt, error) {
	_ = (func(string) interface{} { return nil })(query)
	err := errors.New(`driver-level failure, cannot execute query`)
	return nil, err
}

func (*failingMock) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	_ = (func(context.Context, string) interface{} { return nil })(ctx, query)
	err := errors.New(`driver-level failure, cannot execute query`)
	return nil, err
}

func (*failingMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	_ = (func(string, ...interface{}) interface{} { return nil })(query, args)
	err := errors.New(`driver-level failure, cannot execute query`)
	return nil, err
}

func (*failingMock) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	_ = (func(context.Context, string, ...interface{}) interface{} { return nil })(ctx, query, args)
	err := errors.New(`driver-level failure, cannot execute query`)
	return nil, err
}

func TestNewInsert(t *testing.T) {
	naked := &Insert{
		Table:       tbl,
		Record:      rec,
		recordType:  reflect.TypeOf(rec),
		recordValue: reflect.ValueOf(rec),
	}
	built := NewInsert(tbl, rec)
	if naked.Table != built.Table {
		t.Fatalf(`expected NewInsert builder to return struct with table identical to that of test object`)
	}
	if !reflect.DeepEqual(naked.Record.(candyInsert), built.Record.(candyInsert)) {
		t.Fatalf(`expected NewInsert builder to return struct with record values identical to those of test object`)
	}
}

func TestTokenizeColumnNameTokenType(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `id, candy_name, form_factor, description, manufacturer, weight_grams, ts`
	columnNames := ins.Tokenize(ColumnNameTokenType)
	if expected != columnNames {
		t.Fatalf(`expected "%s", got "%s"`, expected, columnNames)
	}
}

func TestTokenizeQuestionMarkTokenType(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `?, ?, ?, ?, ?, ?, ?`
	bindParams := ins.Tokenize(QuestionMarkTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestTokenizeAtColumnNameTokenType(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `@id, @candy_name, @form_factor, @description, @manufacturer, @weight_grams, @ts`
	bindParams := ins.Tokenize(AtColumnNameTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestTokenizeOrdinalNumberTokenType(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `$1, $2, $3, $4, $5, $6, $7`
	bindParams := ins.Tokenize(OrdinalNumberTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestTokenizeColonTokenType(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `:id, :candy_name, :form_factor, :description, :manufacturer, :weight_grams, :ts`
	bindParams := ins.Tokenize(ColonTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestColumns(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `id, candy_name, form_factor, description, manufacturer, weight_grams, ts`
	columns := ins.Columns()
	if expected != columns {
		t.Fatalf(`expected "%s", got "%s"`, expected, columns)
	}
}

func TestParams(t *testing.T) {
	UseTokenType = QuestionMarkTokenType
	ins := NewInsert(tbl, rec)
	expected := `?, ?, ?, ?, ?, ?, ?`
	params := ins.Params()
	if expected != params {
		t.Fatalf(`expected "%s", got "%s"`, expected, params)
	}
}

func TestSQL(t *testing.T) {
	UseTokenType = OrdinalNumberTokenType
	ins := NewInsert(tbl, rec)
	expected := `INSERT INTO candy (id, candy_name, form_factor, description, manufacturer, weight_grams, ts) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	insertSQL := ins.SQL()
	if expected != insertSQL {
		t.Fatalf(`expected "%s", got "%s"`, expected, insertSQL)
	}
}

func TestArgs(t *testing.T) {
	ins := NewInsert(tbl, rec)
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

func TestInsert(t *testing.T) {
	for tt := range valuesTokenTypes {
		UseTokenType = TokenType(tt)
		ins := NewInsert(tbl, rec)
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
		ins := NewInsert(tbl, rec)
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

func TestInsertError(t *testing.T) {
	fm := &failingMock{}
	UseTokenType = QuestionMarkTokenType
	ins := NewInsert(tbl, rec)
	s := regexp.QuoteMeta(ins.SQL())
	_, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(`failed to construct SQL mock %s`, err)
	}
	mock.ExpectPrepare(s)
	mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = ins.Insert(fm)
	if err == nil {
		t.Fatalf(`expected Insert() to fail, but it succeeded`)
	}
	fmt.Printf(`simulation of driver-level error succeeded: %s`, err)
}

func TestInsertContextError(t *testing.T) {
	fm := &failingMock{}
	UseTokenType = QuestionMarkTokenType
	ins := NewInsert(tbl, rec)
	s := regexp.QuoteMeta(ins.SQL())
	_, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(`failed to construct SQL mock %s`, err)
	}
	mock.ExpectPrepare(s)
	mock.ExpectExec(s).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = ins.InsertContext(context.Background(), fm)
	if err == nil {
		t.Fatalf(`expected InsertContext() to fail, but it succeeded`)
	}
	fmt.Printf(`simulation of driver-level error succeeded: %s`, err)
}
