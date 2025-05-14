package dblite

import (
	"fmt"
	"github.com/franela/goblin"
	"os"
	"testing"
	"time"
)

const sqlModel = `
DROP TABLE IF EXISTS model;
CREATE TABLE IF NOT EXISTS model (
	id            		 INTEGER NOT NULL PRIMARY KEY,
	email         		 TEXT NOT NULL UNIQUE,
	name          		 TEXT DEFAULT '',
	address   			 TEXT DEFAULT '',
	active        		 INTEGER DEFAULT 1
);
`

var dbInstance *Database

type Model struct {
	Id      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  int    `json:"active"`
}

func NewModel(id int64) *Model {
	return &Model{Id: id}
}

func (model *Model) New() *Model {
	return NewModel(-1)
}

func (model *Model) Clone() *Model {
	var o = *model
	return &o
}

func (model *Model) TableName() string {
	return "model"
}

func (model *Model) InsertWithArgs() (bool, int64, error) {
	return Insert(dbInstance.Conn, model, []string{
		`id`, `email`, `name`, `address`,
	}, On{
		On:        "CONFLICT(id) DO UPDATE SET email=?, name=?, address=?",
		Arguments: []any{model.Email, model.Name, model.Address},
	}, "sqlite3")
}

func (model *Model) InsertOnConflictDoNothing() (bool, int64, error) {
	return Insert(dbInstance.Conn, model, []string{
		`id`, `email`, `name`, `address`,
	}, On{
		On: "CONFLICT(id) DO NOTHING",
	}, "sqlite3")
}

func (model *Model) Upsert() (bool, int64, error) {
	var columns = []string{`id`, `email`, `name`, `address`}
	return Insert(dbInstance.Conn, model, columns, On{
		On:            "CONFLICT(id)",
		UpsertColumns: columns[1:],
	}, "sqlite3")
}

func initDB() {
	var dbDIR = "./bin"
	var dbPath = fmt.Sprintf("%v/test.db", dbDIR)
	err := os.MkdirAll(dbDIR, 0755)
	checkError(err)
	dbInstance, err = NewDatabase(dbPath)
	checkError(err)
	_, err = dbInstance.Exec(sqlModel)
	checkError(err)
}

func deInitDB() {
	if dbInstance != nil {
		dbInstance.Close()
	}
}

func TestDBLite(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests Model Insert", func() {
		g.It("model insert", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(1)
			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"

			bln, _, err := m.InsertOnConflictDoNothing()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

			bln, _, err = m.InsertOnConflictDoNothing()
			g.Assert(bln).IsFalse()
			g.Assert(err).IsNil()
		})
		g.It("model count", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var models = []*Model{
				{Id: 1, Email: "email1@db.com", Name: "model1", Address: "123 db street"},
				{Id: 2, Email: "email2@db.com", Name: "model2", Address: "124 db street"},
				{Id: 3, Email: "email3@db.com", Name: "model1", Address: "125 db street"},
				{Id: 4, Email: "email4@db.com", Name: "model4", Address: "126 db street"},
				{Id: 5, Email: "email5@db.com", Name: "model1", Address: "127 db street"},
			}

			for _, model := range models {
				bln, _, err := model.InsertOnConflictDoNothing()
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()
			}
			num, err := Count(dbInstance.Conn, NewModel(-1), `id`, WhereClause{
				Where: `name=?`, Arguments: []any{"model1"},
			})
			g.Assert(err).IsNil()
			g.Assert(num).Equal(int64(3))
			num, err = Count(dbInstance.Conn, NewModel(-1), `id`, WhereClause{
				Where: `name=?`, Arguments: []any{"model4"},
			})
			g.Assert(err).IsNil()
			g.Assert(num).Equal(int64(1))
		})

		g.It("model upsert", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(1)
			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"

			bln, _, err := m.Upsert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

			bln, _, err = m.Upsert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

		})

		g.It("model upsert with sql args", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(1)

			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"

			bln, _, err := m.InsertWithArgs()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

			bln, _, err = m.InsertWithArgs()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()

		})

		g.It("model delete with sql arg", func() {
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
