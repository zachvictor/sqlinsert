# sqlinsert
Generate SQL INSERT statement with bind parameters directly from a Go struct. Define column names in struct tags. Supports bind parameter token types of MySQL, PostgreSQL, Oracle, SingleStore (MemSQL), SQL Server (T-SQL), and their equivalents. Struct values become bind arguments. Works seamlessly with standard [database/sql](https://pkg.go.dev/database/sql) package. Use SQL string and Args slice piecemeal or use `Insert()` or `InsertContext()` with a `sql.Conn`, `sql.DB`, or `sql.Tx` to execute the INSERT statement directly.  
[![Go Reference](https://pkg.go.dev/badge/github.com/zachvictor/sqlinsert.svg)](https://pkg.go.dev/github.com/zachvictor/sqlinsert)
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
type candyInsert struct {
	Id          string    `col:"id"`
	Name        string    `col:"candy_name"`
	FormFactor  string    `col:"form_factor"`
	Description string    `col:"description"`
	Mfr         string    `col:"manufacturer"`
	Weight      float64   `col:"weight_grams"`
	Timestamp   time.Time `col:"ts"`
}

var rec = candyInsert{
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
    VALUES ($1, $2, $3, $4, $5, $6, $7)`)
_, err := stmt.Exec(candyInsert.Id, candyInsert.Name, candyInsert.FormFactor, 
	candyInsert.Description, candyInsert.Mfr, candyInsert.Weight, candyInsert.Timestamp)
```

### After
```go
sqlinsert.UseTokenType = OrdinalNumberTokenType
ins := sqlinsert.NewInsert(`candy`, rec)
_, err := ins.Insert(db)
```