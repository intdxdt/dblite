# `dblite`

A light layer for interacting with `postgres` and `sqlite3`.

## primary key 
there is a difference in adding auto increment for 
`postgres` and `sqlite3`. Add `{PRIMARY KEY}` to be replaced 
based on appropriate driver like this:
```go
var db, err := NewDatabase(driver, uri)
db.SetAutoIncrementPrimaryKey(&model)
```