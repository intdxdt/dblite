package dblite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
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
	case "mariadb", "maria", "mysql":
		driver = "mysql"
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
	var quote string
	switch db.driver {
	case "mysql":
		quote = "`"
	default:
		quote = `"`
	}
	return quote + col + quote
}

func (db *Database) ColumnEqualExcludedAttributes(cols []string) string {
	var columns = make([]string, len(cols))
	for i, col := range cols {
		col = db.QuoteColumn(col)
		switch db.driver {
		case "postgres":
			columns[i] = fmt.Sprintf("%s = EXCLUDED.%s", col, col)
		case "mysql":
			columns[i] = fmt.Sprintf("%s = VALUES(%s)", col, col)
		default:
			columns[i] = fmt.Sprintf("%s = excluded.%s", col, col)
		}
	}
	return strings.Join(columns, ",")
}

func (db *Database) SetClause(col string, index ...int) string {
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

func (db *Database) SetClauses(cols []string, offset ...int) string {
	var placeholders = make([]string, len(cols))
	var n = 0
	if len(offset) > 0 {
		n = offset[0]
	}
	for i, col := range cols {
		placeholders[i] = db.SetClause(col, n+i+1)
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
