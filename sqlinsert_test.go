package sqlinsert

import (
	"reflect"
	"testing"
	"time"
)

type CandyInsert struct {
	Id          string    `col:"id"`
	Name        string    `col:"candy_name"`
	FormFactor  string    `col:"form_factor"`
	Description []string  `col:"description"`
	Mfr         string    `col:"manufacturer"`
	Weight      float64   `col:"weight_grams"`
	Timestamp   time.Time `col:"ts"`
}

var tbl = `candy`

var rec = CandyInsert{
	Id:          `c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
	Name:        `Gougat`,
	FormFactor:  `Package`,
	Description: []string{`tastes`, `like`, `gopher`, `feed`},
	Mfr:         `Gouggle`,
	Weight:      1.16180,
	Timestamp:   time.Time{},
}

func TestNewInsert(t *testing.T) {
	naked := &Insert{
		Table:    tbl,
		Record:   rec,
		RowType:  reflect.TypeOf(rec),
		RowValue: reflect.ValueOf(rec),
	}
	built := NewInsert(tbl, rec)
	if naked.Table != built.Table {
		t.Fatalf(`expected NewInsert builder to return struct with table identical to that of test object`)
	}
	if !reflect.DeepEqual(naked.Record.(CandyInsert), built.Record.(CandyInsert)) {
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
	ins := NewInsert(tbl, rec)
	expected := `?, ?, ?, ?, ?, ?, ?`
	params := ins.Params(QuestionMarkTokenType)
	if expected != params {
		t.Fatalf(`expected "%s", got "%s"`, expected, params)
	}
}

func TestSQL(t *testing.T) {
	ins := NewInsert(tbl, rec)
	expected := `INSERT INTO candy (id, candy_name, form_factor, description, manufacturer, weight_grams, ts) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	insertSQL := ins.SQL(OrdinalNumberTokenType)
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
		[]string{`tastes`, `like`, `gopher`, `feed`},
		`Gouggle`,
		1.1618,
		time.Time{},
	}
	args := ins.Args()
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf(`expected "%s", got "%s"`, expected, args)
	}
}
