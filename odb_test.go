package dblite

import (
	"github.com/franela/goblin"
	"testing"
	"time"
)

const sqlOModel = `
DROP TABLE IF EXISTS omodel;
CREATE TABLE IF NOT EXISTS omodel (
	id            		 INTEGER NOT NULL PRIMARY KEY,
	email         		 TEXT NOT NULL UNIQUE,
	name          		 TEXT DEFAULT '',
	address   			 TEXT DEFAULT '',
	active        		 INTEGER DEFAULT 1
);
`

type OModel struct {
	Id      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  int    `json:"active"`
}

func NewOModel(id int64) *OModel {
	return &OModel{Id: id}
}

func (omodel *OModel) New() *OModel {
	return NewOModel(-1)
}

func (omodel *OModel) Clone() *OModel {
	var o = *omodel
	return &o
}

func (omodel *OModel) TableName() string {
	return "omodel"
}

var dbSource *DatabaseSource

func initPostgresDB() {
	var dbType = "postgres"
	var uri = "user=postgres password=1234 dbname=postgres host=localhost port=5432 sslmode=disable"
	var err error
	dbSource, err = NewDatabaseSource(dbType, uri)
	checkError(err)
	_, err = dbSource.Exec(sqlOModel)
	checkError(err)
}

func deInitPostgresDB() {
	if dbSource != nil {
		dbSource.Close()
	}
}

func TestODB(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Test ODB connection", func() {
		g.It("test insert", func() {
			g.Timeout(1 * time.Hour)
			initPostgresDB()
			defer deInitPostgresDB()

			om := NewOModel(1)
			om.Name = "omodel1"
			om.Email = "omodel@odb.com"
			om.Address = "2221 454 343"
			cols := []string{"id", "email", "name", "address"}
			on := On{
				On: "CONFLICT(id) DO UPDATE SET email=$1, name=$2, address=$3",
			}
			bln, _, err := Insert(dbSource.Conn, om, cols, on, "postgres")
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})

		g.It("test upsert", func() {
			g.Timeout(1 * time.Hour)
			initPostgresDB()
			defer deInitPostgresDB()

			om := NewOModel(1)
			om.Name = "omodel1"
			om.Email = "omodel@odb.com"
			om.Address = "2221 454 343"
			cols := []string{"id", "email", "name", "address"}
			on := On{
				On:            "CONFLICT(id)",
				UpsertColumns: cols[1:],
			}
			bln, _, err := Insert(dbSource.Conn, om, cols, on, "postgres")
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

			bln, _, err = Insert(dbSource.Conn, om, cols, on, "postgres")
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

		})

		g.It("test delete", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(100)
			m.Email = "email100@db.com"
			m.Name = "model100"
			m.Address = "123 db street"
			m.Address = "123 db street"

			bln, _, err := m.InsertOnConflictDoNothing()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

			w := WhereClause{
				Where: `name=?`, Arguments: []any{"model100"},
			}
			n, err := Delete(dbInstance.Conn, m, w)
			g.Assert(n).Equal(int64(1))
			g.Assert(err).IsNil()
		})
	})
}
