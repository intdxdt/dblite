package dblite

import (
	"fmt"
)

func Delete[T ITable[T]](db *Database, model T, wc *Where) (int64, error) {
	var query = fmt.Sprintf(`DELETE FROM %v WHERE %v;`, model.TableName(), wc.clause)

	var res, err = Exec(db.Conn, query, wc.arguments...)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}
