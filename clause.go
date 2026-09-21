package dblite

import (
	"fmt"
)

type FuncGetColumnValues func(inputCols []string) ([]string, []any)

type On struct {
	clause string
	cols   []string
	args   []any
}

type OnOpts func(*On)

func NewOn(onString string, opts ...OnOpts) *On {
	var o = &On{onString, []string{}, []any{}}
	for _, fn := range opts {
		fn(o)
	}
	return o
}

func WithOnColumns(cols []string) OnOpts {
	return func(o *On) {
		o.cols = cols
	}
}

func WithOnArguments(args []any) OnOpts {
	return func(o *On) {
		o.args = args
	}
}

func (on *On) hasOn() bool {
	return len(on.clause) > 0
}

func (on *On) hasColumns() bool {
	return len(on.cols) > 0
}

func (on *On) hasArguments() bool {
	return len(on.args) > 0
}

func (on *On) OnClause(db *Database, getColumnValues FuncGetColumnValues) (string, []any) {
	var onSql string
	var values = make([]any, 0, len(on.args))
	if len(on.cols) > 0 { //do an upsert given upsert columns
		var upsertCols, _ = getColumnValues(on.cols)
		var colPlaceholders = db.ColumnEqualExcludedAttributes(upsertCols)
		onSql = fmt.Sprintf(`%v DO UPDATE SET %v`, on.clause, colPlaceholders)
	} else if len(on.args) > 0 { //on with arguments - maybe not an upsert
		onSql = on.clause
		values = append(values, on.args...)
	} else {
		onSql = on.clause
	}
	return onSql, values
}

type Where struct {
	clause string
	args   []any
}

type WhereOpts func(*Where)

func NewWhere(clause string, opts ...WhereOpts) *Where {
	var o = &Where{clause, []any{}}
	for _, fn := range opts {
		fn(o)
	}
	return o
}

func WithWhereArguments(args []any) WhereOpts {
	return func(o *Where) {
		o.args = args
	}
}
