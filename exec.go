package dblite

import (
	"database/sql"
	"errors"
)

func Exec(conn *sql.DB, query string, args ...any) (sql.Result, error) {
	return conn.Exec(query, args...)
}

// ExecMany return error and rollback error
func ExecMany(conn *sql.DB, query string, records [][]any) (bool, error) {
	tx, err := conn.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		var prepError = tx.Rollback() // rollback if prepare fails
		return false, errors.Join(err, prepError)
	}
	defer stmt.Close()

	var rowCount = int64(0)
	for _, record := range records {
		res, err := stmt.Exec(record...)
		if err != nil {
			return false, errors.Join(err, tx.Rollback())
		}
		n, err := res.RowsAffected()
		if err != nil {
			return false, errors.Join(err, tx.Rollback())
		}
		rowCount = rowCount + n
	}

	err = tx.Commit()
	if err != nil {
		return false, errors.Join(err, tx.Rollback())
	}

	return rowCount == int64(len(records)), nil
}
