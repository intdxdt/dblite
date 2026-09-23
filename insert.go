package dblite

import (
	"fmt"

	ref "github.com/intdxdt/goreflect"
)

type FuncColumnsPlaceholders func(cols string, holders string) (string, string)

func Insert[T ITable[T]](db *Database, model T, insertCols []string, options ...QueryOpt) (bool, error) {
	var opts = NewQueryOption(options...)
	var fields, err = ref.Fields(model)
	if err != nil {
		return false, err
	}

	fields, colRefs, err := ref.FilterFieldReferences(fields, model)
	if err != nil {
		return false, err
	}

	var getColVals = func(inputCols []string) ([]string, []any) {
		var cols = make([]string, 0, len(fields))
		var values = make([]any, 0, len(fields))

		var dict = KeysToMap(inputCols, true)

		for i, field := range fields {
			if !(dict[field]) {
				continue
			}
			cols = append(cols, field)
			values = append(values, colRefs[i])
		}
		return cols, values
	}

	var cols, values = getColVals(insertCols)

	var columns = db.ColumnNames(cols)
	var holders = db.ColumnPlaceholders(cols)
	//:> callback - can modify cols and holders
	columns, holders = opts.funcColumnsPlaceholders(columns, holders)

	var sqlStatement = fmt.Sprintf(`
			INSERT INTO %v(%v) 
			VALUES (%v);`, model.TableName(), columns, holders)

	if opts.hasOn() {
		var onSql, onValues = opts.OnClause(db, getColVals)
		values = append(values, onValues...)
		sqlStatement = fmt.Sprintf(`
			INSERT INTO %v(%v) 
			VALUES (%v)
			ON %v;`, model.TableName(), columns, holders, onSql)
	}

	res, err := Exec(db.Conn, sqlStatement, values...)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func InsertReturning[T ITable[T]](db *Database, model T, insertCols []string, returnColumn string, options ...QueryOpt) (int64, error) {
	var opts = NewQueryOption(options...)
	var fields, err = ref.Fields(model)
	if err != nil {
		return 0, err
	}

	fields, colRefs, err := ref.FilterFieldReferences(fields, model)
	if err != nil {
		return 0, err
	}

	var getColVals = func(inputCols []string) ([]string, []any) {
		var cols = make([]string, 0, len(fields))
		var values = make([]any, 0, len(fields))

		var dict = KeysToMap(inputCols, true)

		for i, field := range fields {
			if !(dict[field]) {
				continue
			}
			cols = append(cols, field)
			values = append(values, colRefs[i])
		}
		return cols, values
	}

	var cols, values = getColVals(insertCols)

	var columns = db.ColumnNames(cols)
	var holders = db.ColumnPlaceholders(cols)
	//:> callback - can modify cols and holders
	columns, holders = opts.funcColumnsPlaceholders(columns, holders)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		RETURNING %s;`, model.TableName(), columns, holders, returnColumn)

	if opts.hasOn() {
		var onSql, onValues = opts.OnClause(db, getColVals)
		values = append(values, onValues...)
		sqlStatement = fmt.Sprintf(`
			INSERT INTO %v(%v) 
			VALUES (%v)
			ON %v
			RETURNING %s;`, model.TableName(), columns, holders, onSql, returnColumn)
	}

	var returnId int64

	err = QueryRow(db.Conn, sqlStatement, values...).Scan(&returnId)
	if err != nil {
		return 0, err
	}

	return returnId, nil
}

func InsertMany[T ITable[T]](db *Database, models []T, insertCols []string, options ...QueryOpt) (bool, error) {
	var opts = NewQueryOption(options...)

	if len(models) == 0 {
		return true, nil
	}

	var getColumnsValues = func(model T) ([]string, []any, error) {
		var fields, err = ref.Fields(model)
		if err != nil {
			return nil, nil, err
		}

		fields, colRefs, err := ref.FilterFieldReferences(fields, model)
		if err != nil {
			return nil, nil, err
		}
		var cols = make([]string, 0, len(fields))
		var values = make([]any, 0, len(fields))

		var dict = KeysToMap(insertCols, true)

		for i, field := range fields {
			if !dict[field] {
				continue
			}
			cols = append(cols, field)
			values = append(values, colRefs[i])
		}
		return cols, values, nil
	}
	var model = models[0]
	var cols, _, err = getColumnsValues(model)
	if err != nil {
		return false, err
	}

	var columns = db.ColumnNames(cols)
	var holders = db.ColumnPlaceholders(cols)
	//:> callback - can modify cols and holders
	columns, holders = opts.funcColumnsPlaceholders(columns, holders)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if opts.hasOnUpsertColumns() {
		panic("only ON clause wth placeholders and apply to all arguments supported")
	}

	if opts.hasOn() {
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, opts.onString())
	}

	var records = make([][]any, 0, len(models))
	for _, model = range models {
		_, values, err := getColumnsValues(model)
		if err != nil {
			return false, err
		}
		if opts.hasOn() {
			for _, v := range opts.on.args {
				values = append(values, v)
			}
		}
		records = append(records, values)
	}

	return ExecMany(db.Conn, sqlStatement, records)
}
