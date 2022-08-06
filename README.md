# sqlinsert
Generate a SQL INSERT statement with bind parameters directly from a Go struct.
[![Go Reference](https://pkg.go.dev/badge/github.com/zachvictor/sqlinsert.svg)](https://pkg.go.dev/github.com/zachvictor/sqlinsert)

## Features
* Define column names in struct tags.
* Struct values become bind arguments.
* Use SQL outputs and Args slice piecemeal. Or, use `Insert()`/`InsertContext()` with a `sql.Conn`, `sql.DB`, or 
`sql.Tx` to execute the INSERT statement directly.
* Works seamlessly with Go standard library [database/sql](https://pkg.go.dev/database/sql) package. 
* Supports bind parameter token types of MySQL, PostgreSQL, Oracle, SingleStore (MemSQL), SQL Server (T-SQL), and their 
equivalents.
* Supports customized struct tags and token types.
* Supports Go 1.8 to 1.19.
* Test coverage: 100% files, 97.5% statements. Tested on Go 1.15, 1.17, and 1.18.

## Example
### Given
```sql
CREATE TABLE candy (
    id           CHAR(36) NOT NULL
    candy_name   VARCHAR(255) NOT NULL
    form_factor  VARCHAR(255) NOT NULL
    description  VARCHAR(255) NOT NULL
    manufacturer VARCHAR(255) NOT NULL
    weight_grams DECIMAL(9, 3) NOT NULL
    ts DATETIME  NOT NULL
)
```

```go
type CandyInsert struct {
    Id          string    `col:"id"`
    Name        string    `col:"candy_name"`
    FormFactor  string    `col:"form_factor"`
    Description string    `col:"description"`
    Mfr         string    `col:"manufacturer"`
    Weight      float64   `col:"weight_grams"`
    Timestamp   time.Time `col:"ts"`
}

var rec = CandyInsert{
    Id:          `c0600afd-78a7-4a1a-87c5-1bc48cafd14e`,
    Name:        `Gougat`,
    FormFactor:  `Package`,
    Description: `tastes like gopher feed`,
    Mfr:         `Gouggle`,
    Weight:      1.16180,
    Timestamp:   time.Time{},
}
```

### Before
```go
stmt, _ := db.Prepare(`INSERT INTO candy
    (id, candy_name, form_factor, description, manufacturer, weight_grams, ts)
    VALUES (?, ?, ?, ?, ?, ?, ?)`)
_, err := stmt.Exec(candyInsert.Id, candyInsert.Name, candyInsert.FormFactor,
	candyInsert.Description, candyInsert.Mfr, candyInsert.Weight, candyInsert.Timestamp)
```

### After
```go
ins := sqlinsert.Insert{`candy`, &rec}
_, err := ins.Insert(db)
```

## This is not an ORM

### Hide nothing
Unlike ORMs, `sqlinsert` does **not** create an abstraction layer over SQL relations, nor does it restructure SQL
functions.
The aim is to keep it simple and hide nothing.
`sqlinsert` is fundamentally a helper for [database/sql](https://pkg.go.dev/database/sql).
It simply maps struct fields to INSERT elements:
* struct tags
=> SQL columns and tokens `string`
=> [Prepare](https://pkg.go.dev/database/sql@go1.17#DB.Prepare) `query string`
* struct values
=> bind args `[]interface{}`
=> [Exec](https://pkg.go.dev/database/sql@go1.17#Stmt.Exec) `args ...interface{}`
([Go 1.18](https://pkg.go.dev/database/sql@go1.18#DB.ExecContext)+ `args ...any`)

All aspects of SQL INSERT remain in your control:
* *I just want the column names for my SQL.* `Insert.Columns()`
* *I just want the parameter-tokens for my SQL.* `Insert.Params()`
* *I just want the bind args for my Exec() call.* `Insert.Args()`
* *I just want a simple, one-for-one wrapper.* `Insert.Insert()`

### Let SQL be great
SQL’s INSERT is already as close to functionally pure as possible. Why would we change that? Its simplicity and
directness are its power.

### Let database/sql be great
Some database vendors support collection types for bind parameters, some don’t.
Some database drivers support slices for bind args, some don’t.
The complexity of this reality is met admirably by [database/sql](https://pkg.go.dev/database/sql)
with the _necessary_ amount of flexibility and abstraction:
*flexibility* in open-ended SQL;
*abstraction* in the variadic `args ...interface{}` for bind args.
In this way, [database/sql](https://pkg.go.dev/database/sql) respects INSERT’s power,
hiding nothing even as it tolerates the vagaries of bind-parameter handling among database vendors and drivers.

### Let Go be great
Go structs support ordered fields, strong types, and field metadata via [tags](https://go.dev/ref/spec#Tag) and
[reflection](https://pkg.go.dev/reflect#StructTag).
In these respects, the Go struct can encapsulate the information of a SQL INSERT-row perfectly and completely.
`sqlinsert` uses these features of Go structs to makes your SQL INSERT experience more Go-idiomatic.

## Limitations
`sqlinsert` is for simple binding. It does not support SQL operations in the `VALUES` clause.
If you require, say—
```sql
INSERT INTO foo (bar, baz, oof) VALUES (some_function(?), REPLACE(?, 'oink', 'moo'), ? + ?);
```
—then you can use `sqlinsert.Insert` methods piecemeal.
For example, use `Insert.Columns` to build the column list for `Prepare`
and `Insert.Args` to marshal the args for `Exec`/`ExecContext`.
