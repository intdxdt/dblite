package dblite

import (
	"fmt"
)

type On struct {
	On            string
	UpsertColumns []string
	Arguments     []any
}

func (on *On) hasOn() bool {
	return len(on.On) > 0
}

func (on *On) OnClause(db *Database, getColumnValues func(inputCols []string) ([]string, []any)) (string, []any) {
	var onSql string
	var values = make([]any, 0, len(on.Arguments))
	if len(on.UpsertColumns) > 0 { //do an upsert given upsert columns
		var upsertCols, _ = getColumnValues(on.UpsertColumns)
		var colPlaceholders = db.ColumnEqualExcludedAttributes(upsertCols)
		onSql = fmt.Sprintf(`%v DO UPDATE SET %v`, on.On, colPlaceholders)
	} else if len(on.Arguments) > 0 { //on with arguments - maybe not an upsert
		onSql = on.On
		values = append(values, on.Arguments...)
	} else {
		onSql = on.On
	}
	return onSql, values
}

type WhereClause struct {
	Where     string
	Arguments []any
}
