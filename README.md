# sqlinsert
Generate SQL INSERT statement with bind parameters directly from a Go struct. Define column names in struct tags. Supports bind parameter token types of MySQL, PostgreSQL, Oracle, SingleStore (MemSQL), SQL Server (T-SQL), and their equivalents. Struct values become bind arguments. Works seamlessly with standard [database/sql](https://pkg.go.dev/database/sql) package.
