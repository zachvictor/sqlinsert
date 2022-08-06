package sqlinsert

import "time"

var valuesTokenTypes = []TokenType{
	QuestionMarkTokenType,
	AtColumnNameTokenType,
	OrdinalNumberTokenType,
	ColonTokenType,
}

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

var recValue = candyInsert{
	Id:          `c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
	Name:        `Gougat`,
	FormFactor:  `Package`,
	Description: `tastes like gopher feed`,
	Mfr:         `Gouggle`,
	Weight:      1.16180,
	Timestamp:   time.Time{},
}

var recPointer = &candyInsert{
	Id:          `c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
	Name:        `Gougat`,
	FormFactor:  `Package`,
	Description: `tastes like gopher feed`,
	Mfr:         `Gouggle`,
	Weight:      1.16180,
	Timestamp:   time.Time{},
}

var fiveRecsValues = []candyInsert{
	{`a`, `a`, `a`, `a`, `a`, 1.1, time.Time{}},
	{`b`, `b`, `b`, `b`, `b`, 2.1, time.Time{}},
	{`c`, `c`, `c`, `c`, `c`, 3.1, time.Time{}},
	{`d`, `d`, `d`, `d`, `d`, 4.1, time.Time{}},
	{`e`, `e`, `e`, `e`, `e`, 5.1, time.Time{}},
}

var fiveRecsPointers = []*candyInsert{
	{`a`, `a`, `a`, `a`, `a`, 1.1, time.Time{}},
	{`b`, `b`, `b`, `b`, `b`, 2.1, time.Time{}},
	{`c`, `c`, `c`, `c`, `c`, 3.1, time.Time{}},
	{`d`, `d`, `d`, `d`, `d`, 4.1, time.Time{}},
	{`e`, `e`, `e`, `e`, `e`, 5.1, time.Time{}},
}
