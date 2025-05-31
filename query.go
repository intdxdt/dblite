package dblite

import (
	"database/sql"
	"fmt"

	ref "github.com/intdxdt/goreflect"
)

func Query(db *Database, query string, args ...any) (*sql.Rows, error) {
	return db.Conn.Query(query, args...)
}

func QueryModel[T ITable[T]](db *Database, model T, where ...WhereClause) (T, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return model.New(), err
	}
	return QueryModelByColumnNames(db, model, fields, where...)
}

func QueryModelByColumnNames[T ITable[T]](db *Database, model T, fieldNames []string, where ...WhereClause) (T, error) {
	var tableName = model.TableName()
	var cols, colRefs, err = ref.FilterFieldReferences(fieldNames, model)
	if err != nil {
		return model, err
	}
	var fields = db.ColumnNames(cols)

	var args = make([]any, 0)
	var sqlStatement = fmt.Sprintf("SELECT %v FROM %v LIMIT 1;", fields, tableName)
	if len(where) > 0 {
		var wc = where[0]
		args = wc.Arguments
		sqlStatement = fmt.Sprintf("SELECT %v FROM %v WHERE %v LIMIT 1;", fields, tableName, wc.Where)
	}

	rows, err := Query(db, sqlStatement, args...)
	if err != nil {
		return model, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(colRefs...)
		if err != nil {
			return model, err
		}
		break
	}

	if rows.Err() != nil {
		return model, rows.Err()
	}
	return model, nil
}

func QueryModels[T ITable[T]](db *Database, model T, where ...WhereClause) ([]T, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return []T{}, err
	}
	return QueriesByColumnNames(db, model, fields, where...)
}

func QueriesByColumnNames[T ITable[T]](db *Database, model T, fieldNames []string, where ...WhereClause) ([]T, error) {
	var results = make([]T, 0)
	var tableName = model.TableName()
	var cols, colRefs, err = ref.FilterFieldReferences(fieldNames, model)
	if err != nil {
		return nil, err
	}
	var fields = db.ColumnNames(cols)

	var args = make([]any, 0)
	var sqlStatement = fmt.Sprintf("SELECT %v FROM %v;", fields, tableName)
	if len(where) > 0 {
		var wc = where[0]
		args = wc.Arguments
		sqlStatement = fmt.Sprintf("SELECT %v FROM %v WHERE %v;", fields, tableName, wc.Where)
	}

	rows, err := Query(db, sqlStatement, args...)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(colRefs...)
		if err != nil {
			return results, err
		}
		results = append(results, model.Clone())
	}

	if rows.Err() != nil {
		return results, rows.Err()
	}
	return results, nil
}
