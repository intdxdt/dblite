package dblite

import (
	"database/sql"
	"fmt"
)

func Insert[T ITable[T]](conn *sql.DB, model T, insertCols []string, on On) (bool, error) {
	var fields, colRefs, err = model.FilterFieldReferences(model.Fields())
	if err != nil {
		return false, err
	}
	var cols = make([]string, 0, len(fields))
	var values = make([]any, 0, len(fields))

	var dict = KeysToMap(insertCols, true)

	for i, field := range fields {
		if !(dict[field]) {
			continue
		}
		cols = append(cols, field)
		values = append(values, colRefs[i])
	}

	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if len(on.On) > 0 {
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, on.On)
		for _, v := range on.Arguments {
			values = append(values, v)
		}
	}

	res, err := conn.Exec(sqlStatement, values...)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func InsertMany[T ITable[T]](database *Database, models []T, insertCols []string, on On) error {
	if len(models) == 0 {
		return nil
	}

	var getColumnsValues = func(model T) ([]string, []any, error) {
		var fields, colRefs, err = model.FilterFieldReferences(model.Fields())
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
		return err
	}

	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if len(on.On) > 0 {
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, on.On)
	}

	var records = make([][]any, 0, len(models))
	for _, model = range models {
		_, values, err := getColumnsValues(model)
		if err != nil {
			return err
		}
		if len(on.On) > 0 {
			for _, v := range on.Arguments {
				values = append(values, v)
			}
		}
		records = append(records, values)
	}

	return database.ExecMany(sqlStatement, records)
}
