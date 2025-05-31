package dblite

import (
	"fmt"

	ref "github.com/intdxdt/goreflect"
)

func Update[T ITable[T]](db *Database, model T, updateCols []string, wc WhereClause) (bool, error) {
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
	for _, arg := range wc.Arguments {
		values = append(values, arg)
	}

	var query = fmt.Sprintf(
		`UPDATE %v SET %v WHERE %v;`,
		model.TableName(), holders, wc.Where)

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
