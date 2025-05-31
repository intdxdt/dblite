package dblite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	driver string
	uri    string
	Conn   *sql.DB
}

func NewDatabase(driver, uri string) (*Database, error) {
	driver = strings.ToLower(driver)
	switch driver {
	case "lite", "lite3", "sqlite", "sqlite3":
		driver = "sqlite3"
	case "pgx", "pq", "postgres":
		driver = "postgres"
	default:
		return nil, errors.New("unsupported database driver: " + driver)
	}
	conn, err := sql.Open(driver, uri)
	if err != nil {
		return nil, err
	}
	return &Database{
		driver: driver,
		uri:    uri,
		Conn:   conn,
	}, nil
}

func (db *Database) Close() {
	if db.Conn != nil {
		_ = db.Conn.Close()
	}
}

func (db *Database) ColumnNames(cols []string) string {
	var quoted = make([]string, len(cols))
	for i, col := range cols {
		quoted[i] = db.QuoteColumn(col)

	}
	return strings.Join(quoted, ",")
}

func (db *Database) QuoteColumn(col string) string {
	return fmt.Sprintf(`"%v"`, col)
}

func (db *Database) ColumnEqualExcludedAttributes(cols []string) string {
	var columns = make([]string, len(cols))
	for i, col := range cols {
		col = db.QuoteColumn(col)
		switch db.driver {
		case "postgres":
			columns[i] = fmt.Sprintf("%s = EXCLUDED.%s", col, col)
		default:
			columns[i] = fmt.Sprintf("%s = excluded.%s", col, col)
		}
	}
	return strings.Join(columns, ",")
}

func (db *Database) WhereParam(col string, ops string, index ...int) string {
	var placeholder string
	switch db.driver {
	case "postgres":
		var idx = 1
		if len(index) > 0 {
			idx = index[0]
		}
		placeholder = fmt.Sprintf(`%v%v$%d`, db.QuoteColumn(col), ops, idx)
	default:
		placeholder = fmt.Sprintf(`%v%v?`, db.QuoteColumn(col), ops)
	}
	return placeholder
}

func (db *Database) SetParam(col string, index ...int) string {
	var placeholder string
	switch db.driver {
	case "postgres":
		var idx = 1
		if len(index) > 0 {
			idx = index[0]
		}
		placeholder = fmt.Sprintf(`%v=$%d`, db.QuoteColumn(col), idx)
	default:
		placeholder = fmt.Sprintf(`%v=?`, db.QuoteColumn(col))
	}
	return placeholder
}

func (db *Database) SetParams(cols []string, offset ...int) string {
	var placeholders = make([]string, len(cols))
	var n = 0
	if len(offset) > 0 {
		n = offset[0]
	}
	for i, col := range cols {
		placeholders[i] = db.SetParam(col, n+i+1)
	}
	return strings.Join(placeholders, `,`)
}

func (db *Database) ColumnPlaceholders(cols []string) string {
	switch db.driver {
	case "postgres":
		var placeholders = make([]string, len(cols))
		for i := range cols {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		return strings.Join(placeholders, ",")
	default:
		return strings.TrimRight(strings.Repeat("?,", len(cols)), ",")
	}
}
