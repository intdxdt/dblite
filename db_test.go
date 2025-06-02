package dblite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/joho/godotenv"
)

func sqlModel() string {
	return `
	DROP TABLE IF EXISTS model;
	CREATE TABLE IF NOT EXISTS model (
		id       {PRIMARY KEY},
		email    TEXT NOT NULL UNIQUE,
		name     TEXT DEFAULT '',
		address  TEXT DEFAULT '',
		active   INTEGER DEFAULT 1
	);`
}

func init() {
	var err = godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env", err)
	}
}

var testDrivers = []string{"sqlite3", "postgres"}

type TestModel struct {
	Id      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  int    `json:"active"`
}

func generateData(n int) []*TestModel {
	var models = make([]*TestModel, n)

	for i := 0; i < n; i++ {
		models[i] = &TestModel{
			Id:      int64(i + 1),
			Email:   fmt.Sprintf("user%d@example.com", i),
			Name:    fmt.Sprintf("User %d", i),
			Address: fmt.Sprintf("%d db street, gh", i),
			Active:  i % 2,
		}
	}
	return models
}

func NewModel(id int64) *TestModel {
	return &TestModel{Id: id}
}

func (model *TestModel) New() *TestModel {
	return NewModel(-1)
}

func (model *TestModel) Clone() *TestModel {
	var o = *model
	return &o
}

func (model *TestModel) TableName() string {
	return "model"
}

func initDB(driver string) *Database {

	switch driver {
	case "sqlite3":
		var uri = os.Getenv("SQLITE_URI")
		var dbDIR = filepath.Dir(uri)
		var dbPath = fmt.Sprintf("%v/test.db", dbDIR)
		checkError(os.MkdirAll(dbDIR, 0755))

		db, err := NewDatabase(driver, dbPath)
		checkError(err)

		var model = sqlModel()
		db.SetAutoIncrementPrimaryKey(&model)
		_, err = Exec(db.Conn, model)
		checkError(err)
		return db

	case "postgres":
		var uri = os.Getenv("POSTGRES_URI")
		var db, err = NewDatabase(driver, uri)
		checkError(err)

		var model = sqlModel()
		db.SetAutoIncrementPrimaryKey(&model)
		_, err = Exec(db.Conn, model)
		checkError(err)
		return db
	default:
		log.Fatalln("untested driver", driver)
	}
	return nil
}

func deInitDB(db *Database) {
	if db != nil {
		db.Close()
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func TestDB(t *testing.T) {
	var g = goblin.Goblin(t)
	g.Describe("Test Delete", func() {

		g.It("delete: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var name, err = TableNameFromCreateSql(sqlModel())
				g.Assert(name).Equal("model")
				g.Assert(err).IsNil()

				columns, err := ColumnsByExclusion(NewModel(-1), []string{`id`, `active`})
				g.Assert(columns).Equal([]string{"email", "name", "address"})
				g.Assert(len(columns)).Equal(3)
				g.Assert(err).IsNil()
				deInitDB(db)
			}

		})
	})
}
