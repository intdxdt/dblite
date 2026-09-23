package dblite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	ref "github.com/intdxdt/goreflect"
)

var ErrNoRecord = errors.New("record not found")

func Query(db *Database, query string, args ...any) (*sql.Rows, error) {
	return db.Conn.Query(query, args...)
}

func QueryRow(conn *sql.DB, query string, args ...any) *sql.Row {
	return QueryRowContext(conn, context.Background(), query, args...)
}

func QueryRowContext(conn *sql.DB, ctx context.Context, query string, args ...any) *sql.Row {
	return conn.QueryRowContext(ctx, query, args...)
}

func QueryModel[T ITable[T]](db *Database, model T, options ...QueryOpt) (T, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return model.New(), err
	}
	return QueryModelByColumnNames(db, model, fields, options...)
}

func QueryModelByColumnNames[T ITable[T]](db *Database, model T, fieldNames []string, options ...QueryOpt) (T, error) {
	var opts = NewQueryOption(options...)
	var tableName = model.TableName()
	var cols, colRefs, err = ref.FilterFieldReferences(fieldNames, model)
	if err != nil {
		return model.New(), err
	}

	var columns = db.ColumnNames(cols)
	//:> callback - can modify cols only
	columns, _ = opts.funcColumnsPlaceholders(columns, "")

	var args = make([]any, 0)
	var sqlStatement = fmt.Sprintf("SELECT %v FROM %v LIMIT 1;", columns, tableName)
	if opts.hasWhereClause() {
		args = opts.whereArguments()
		sqlStatement = fmt.Sprintf("SELECT %v FROM %v WHERE %v LIMIT 1;", columns, tableName, opts.whereString())
	}

	rows, err := Query(db, sqlStatement, args...)
	if err != nil {
		return model.New(), err
	}
	defer rows.Close()

	var scanned = false
	for rows.Next() {
		err = rows.Scan(colRefs...)
		if err != nil {
			return model.New(), err
		}
		scanned = true
		break
	}

	if !scanned {
		return model.New(), ErrNoRecord
	}

	if rows.Err() != nil {
		return model.New(), rows.Err()
	}
	return model, nil
}

func QueryModels[T ITable[T]](db *Database, model T, options ...QueryOpt) ([]T, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return []T{}, err
	}
	return QueriesByColumnNames(db, model, fields, options...)
}

func QueriesByColumnNames[T ITable[T]](db *Database, model T, fieldNames []string, options ...QueryOpt) ([]T, error) {
	var opts = NewQueryOption(options...)
	var results = make([]T, 0)
	var tableName = model.TableName()
	var cols, colRefs, err = ref.FilterFieldReferences(fieldNames, model)
	if err != nil {
		return nil, err
	}

	var columns = db.ColumnNames(cols)
	//:> callback - can modify cols only
	columns, _ = opts.funcColumnsPlaceholders(columns, "")

	var args = make([]any, 0)
	var sqlStatement = fmt.Sprintf("SELECT %v FROM %v;", columns, tableName)
	if opts.hasWhereClause() {
		args = opts.whereArguments()
		sqlStatement = fmt.Sprintf("SELECT %v FROM %v WHERE %v;", columns, tableName, opts.whereString())
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
