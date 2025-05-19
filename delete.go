package dblite

import (
	"fmt"
)

func Delete[T ITable[T]](db *Database, model T, wc WhereClause) (int64, error) {
	var query = fmt.Sprintf(
		`DELETE FROM %v WHERE %v;`, model.TableName(), wc.Where)

	var res, err = Exec(db.Conn, query, wc.Arguments...)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}
