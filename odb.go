package dblite

import (
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseSource struct {
	DBType string
	URI    string
	Conn   *sql.DB
}

func NewDatabaseSource(DBType string, URI string) (*DatabaseSource, error) {
	dbConn, err := sql.Open(DBType, URI)
	if err != nil {
		return nil, err
	}
	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}
	return &DatabaseSource{DBType: DBType, URI: URI, Conn: dbConn}, nil
}

func (ds *DatabaseSource) Close() {
	if ds.Conn != nil {
		checkError(ds.Conn.Close())
	}
}

func (ds *DatabaseSource) Exec(query string, args ...any) (sql.Result, error) {
	return Exec(ds.Conn, query, args...)
}

func (ds *DatabaseSource) ExecMany(query string, records [][]any) (error, error) {
	return ExecMany(ds.Conn, query, records)
}
func (ds *DatabaseSource) Query(query string, args ...any) (*sql.Rows, error) {
	return ds.Conn.Query(query, args...)
}
