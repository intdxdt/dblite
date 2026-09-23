package dblite

import (
	"errors"
	"fmt"

	ref "github.com/intdxdt/goreflect"
)

func Update[T ITable[T]](db *Database, model T, updateCols []string, options ...QueryOpt) (bool, error) {
	var opts = NewQueryOption(options...)
	if !opts.hasWhereClause() {
		return false, errors.New("no WHERE clause provided")
	}
	//var wc  = opts.where;
	var fields, err = ref.Fields(model)
	if err != nil {
		return false, err
	}

	fields, colRefs, err := ref.FilterFieldReferences(fields, model)
	if err != nil {
		return false, err
	}

	var cols = make([]string, 0, len(fields))
	var values = make([]any, 0, len(fields))

	var dict = KeysToMap(updateCols, true)

	for i, field := range fields {
		if dict[field] {
			cols = append(cols, field)
			values = append(values, colRefs[i])
		}
	}

	var holders = db.SetParams(cols)

	for _, arg := range opts.whereArguments() {
		values = append(values, arg)
	}

	//:> callback - modify only holders only based on cols
	_, holders = opts.funcColumnsPlaceholders(db.ColumnNames(cols), holders)

	var query = fmt.Sprintf(
		`UPDATE %v SET %v WHERE %v;`, model.TableName(), holders, opts.whereString(),
	)

	res, err := Exec(db.Conn, query, values...)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
