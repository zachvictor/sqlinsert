package sqlinsert

import (
	"reflect"
	"testing"
)

func TestTokenizeColumnNameTokenType(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := `(id,candy_name,form_factor,description,manufacturer,weight_grams,ts)`
	columnNames := Tokenize(reflect.TypeOf(ins.Data), ColumnNameTokenType)
	if expected != columnNames {
		t.Fatalf(`expected "%s", got "%s"`, expected, columnNames)
	}
}

func TestTokenizeQuestionMarkTokenType(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := `(?,?,?,?,?,?,?)`
	bindParams := Tokenize(reflect.TypeOf(ins.Data), QuestionMarkTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestTokenizeAtColumnNameTokenType(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := `(@id,@candy_name,@form_factor,@description,@manufacturer,@weight_grams,@ts)`
	bindParams := Tokenize(reflect.TypeOf(ins.Data), AtColumnNameTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestTokenizeOrdinalNumberTokenType(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := `($1,$2,$3,$4,$5,$6,$7)`
	bindParams := Tokenize(reflect.TypeOf(ins.Data), OrdinalNumberTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}

func TestTokenizeColonTokenType(t *testing.T) {
	ins := Insert{tbl, recValue}
	expected := `(:id,:candy_name,:form_factor,:description,:manufacturer,:weight_grams,:ts)`
	bindParams := Tokenize(reflect.TypeOf(ins.Data), ColonTokenType)
	if expected != bindParams {
		t.Fatalf(`expected "%s", got "%s"`, expected, bindParams)
	}
}
