package sqlinsert

import (
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

// Tokenize translates struct fields into the tokens of SQL column or value expressions as a comma-separated list
// enclosed in parentheses.
func Tokenize(recordType reflect.Type, tokenType TokenType) string {
	var b strings.Builder
	b.WriteString(`(`)
	for i := 0; i < recordType.NumField(); i++ {
		switch tokenType {
		case ColumnNameTokenType:
			b.WriteString(recordType.Field(i).Tag.Get(UseStructTag))
		case QuestionMarkTokenType:
			_, _ = fmt.Fprint(&b, `?`)
		case AtColumnNameTokenType:
			_, _ = fmt.Fprintf(&b, `@%s`, recordType.Field(i).Tag.Get(UseStructTag))
		case OrdinalNumberTokenType:
			_, _ = fmt.Fprintf(&b, `$%d`, i+1)
		case ColonTokenType:
			_, _ = fmt.Fprintf(&b, `:%s`, recordType.Field(i).Tag.Get(UseStructTag))
		}
		if i < recordType.NumField()-1 {
			b.WriteString(`,`)
		}
	}
	b.WriteString(`)`)
	return b.String()
}
