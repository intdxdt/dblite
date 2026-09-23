# dblite

A light layer for interacting with `postgres` and `sqlite3`.
get the latest release 
```bash
go get -u github.com/intdxdt/dblite 
```

# table

```go 
var sqlTable = `
	DROP TABLE IF EXISTS model;
	CREATE TABLE IF NOT EXISTS model (
		id       {PRIMARY KEY},
		email    TEXT NOT NULL UNIQUE,
		name     TEXT DEFAULT '',
		address  TEXT DEFAULT '',
		active   INTEGER DEFAULT 1
	);`
```

## auto increment primary key

there is a difference in adding auto increment for `postgres` and `sqlite3`.
Add `{PRIMARY KEY}` to be replaced based on appropriate driver like this:

```go
db, err := NewDatabase(driver, uri)
db.SetAutoIncrementPrimaryKey(&sqlTable)
```

# model

Model struct should have `json` fields same as table fields

```go
type Model struct {
    Id      int64  `json:"id"`
    Email   string `json:"email"`
    Name    string `json:"name"`
    Address string `json:"address"`
    Active  int    `json:"active"`
}
func NewModel(id int64) *TestModel {
	return &TestModel{Id: id}
}

func (model *TestModel) New() *TestModel {
	return NewModel(-1)
}

func (model *TestModel) Clone() *TestModel {
	return new(*model)
}

func (model *TestModel) TableName() string {
	return db.TableNameFromCreateSql(sqlTable)
}
```

Model should satisfy `ITable` interface:
```go
type ITable[T any] interface {
	New() T
	Clone() T
	TableName() string
}
```

# insert 
Insert new item 
```go 
var m = &Model{
    Id:      1,
    Email:   "email@db.com",
    Name:    "model",
    Address: "123 db street",
}

var bln, err = dblite.Insert(db, m, []string{`id`, `email`, `name`, `address`},
    WithOn(NewOn("CONFLICT(id) DO NOTHING")))
```


## insert many 
```go 
var bln, err = dblite.InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
    WithOn(NewOn("CONFLICT(id) DO NOTHING")))
```

### insert returning id 
```go 
var insertedId, err = dblite.InsertReturning(db, model, []string{`id`, `email`, `name`, `address`}, "id",
    WithOn(NewOn("CONFLICT(id) DO NOTHING"), ))
```

# query 
One : 
```go 
res, err := QueryModel(db, NewModel(-1), WithWhere(NewWhere(
    db.SetParam("id"), WithWhereArguments([]any{73}),
)))
```

Many : 
```go 
res, err := dblite.QueryModels(db, NewModel(-1), WithWhere(NewWhere(`"active"=1`)))
```

# update 
```go 
bln, err := dblite.Update(db, model, cols, WithWhere(NewWhere(
    db.SetParam("id", len(cols)+1), WithWhereArguments([]any{model.Id}),
)))
```

# delete 
```go 
num, err := dblite.Delete(db, NewModel(-1), NewWhere(
    db.WhereParam("active", "="), WithWhereArguments([]any{1}),
))
```

# count 
```go 
num, err := dblite.Count(db, NewModel(-1), `id`, &Where{
    clause: db.SetParam("name"), args: []any{"model1"},
})
```


