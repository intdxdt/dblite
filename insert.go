package dblite

import (
	"database/sql"
	"fmt"
	ref "github.com/intdxdt/goreflect"
)

func Insert[T ITable[T]](conn *sql.DB, model T, insertCols []string, on On, dbType string) (bool, int64, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return false, -1, err
	}

	fields, colRefs, err := ref.FilterFieldReferences(fields, model)
	if err != nil {
		return false, -1, err
	}

	var getColsVals = func(inputCols []string) ([]string, []any) {
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

	var cols, values = getColsVals(insertCols)
	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols, dbType)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if len(on.On) > 0 {
		var sqlOn string
		if len(on.UpsertColumns) > 0 { //do an upsert given upsert columns
			var upsertCols, upsertValues = getColsVals(on.UpsertColumns)
			var colPlaceholders = ColumnEqualParamAttributes(upsertCols, dbType)
			sqlOn = fmt.Sprintf(`%v DO UPDATE SET %v`, on.On, colPlaceholders)
			values = append(values, upsertValues...)
		} else if len(on.Arguments) > 0 { //on with arguments - maybe not an upsert
			sqlOn = on.On
			values = append(values, on.Arguments...)
		} else {
			sqlOn = on.On
		}
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, sqlOn)
	}

	res, err := Exec(conn, sqlStatement, values...)
	if err != nil {
		return false, -1, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, -1, err
	}
	var insertId int64
	if dbType == "postgres" {
		insertId = 0
	} else {
		insertId, err = res.LastInsertId()
		if err != nil {
			return false, -1, err
		}
	}
	return count == 1, insertId, nil
}

func InsertMany[T ITable[T]](conn *sql.DB, models []T, insertCols []string, on On, dbType string) (error, error) {
	if len(models) == 0 {
		return nil, nil
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
		return err, nil
	}

	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols, dbType)

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
			return err, nil
		}
		if len(on.On) > 0 {
			for _, v := range on.Arguments {
				values = append(values, v)
			}
		}
		records = append(records, values)
	}

	return ExecMany(conn, sqlStatement, records)
}
