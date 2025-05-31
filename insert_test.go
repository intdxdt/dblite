package dblite

import (
	"fmt"
	"testing"
	"time"

	"github.com/franela/goblin"
)

func TestInsert(t *testing.T) {
	var g = goblin.Goblin(t)
	g.Describe("Test Insert", func() {
		g.It("insert: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var m = &TestModel{
					Id:      1,
					Email:   "email@db.com",
					Name:    "model",
					Address: "123 db street",
				}

				var bln, err = Insert(db, m, []string{`id`, `email`, `name`, `address`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				bln, err = Insert(db, m, []string{`id`, `email`, `name`, `address`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsFalse()
				g.Assert(err).IsNil()

				deInitDB(db)
			}

		})

		g.It("insert many: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var data = generateData(1024)
				var bln, err = InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				bln, err = InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsFalse()
				g.Assert(err).IsNil()

				var columns = []string{`id`, `email`, `name`, `address`, `active`}
				var setCol = db.SetParam(`active`, len(columns)+1)
				var onClause = fmt.Sprintf("CONFLICT(id) DO UPDATE SET %v", setCol)
				bln, err = InsertMany(db, data, columns, On{
					On:        onClause,
					Arguments: []any{1},
				})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				deInitDB(db)
			}

		})
	})
}

func TestUpsert(t *testing.T) {
	var g = goblin.Goblin(t)

	g.Describe("Test Upsert", func() {
		g.It("upsert", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var m = &TestModel{
					Id:      1,
					Email:   "email@db.com",
					Name:    "model",
					Address: "123 db street",
				}

				var columns = []string{`id`, `email`, `name`, `address`}

				var bln, err = Insert(db, m, columns, On{
					On:            "CONFLICT(id)",
					UpsertColumns: columns[1:],
				})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				bln, err = Insert(db, m, columns, On{
					On:            "CONFLICT(id)",
					UpsertColumns: columns[1:],
				})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				deInitDB(db)
			}
		})

		g.It("upsert with sql args", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)
				var m = &TestModel{
					Id:      1,
					Email:   "email@db.com",
					Name:    "model",
					Address: "123 db street",
				}

				var columns = []string{`id`, `email`, `name`, `address`}
				var setColumns = db.SetParams([]string{`email`, `name`, `address`}, len(columns))
				var bln, err = Insert(db, m, columns, On{
					On:        fmt.Sprintf("CONFLICT(id) DO UPDATE SET %v", setColumns),
					Arguments: []any{m.Email, m.Name, m.Address},
				})

				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				columns = []string{`id`, `email`, `name`, `address`}
				setColumns = db.SetParams([]string{`email`, `name`, `address`}, len(columns))
				bln, err = Insert(db, m, columns, On{
					On:        fmt.Sprintf("CONFLICT(id) DO UPDATE SET %v", setColumns),
					Arguments: []any{m.Email, m.Name, m.Address},
				})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				deInitDB(db)
			}
		})

		g.It("upsert just on", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)
				var m = &TestModel{
					Id:      1,
					Email:   "email@db.com",
					Name:    "model",
					Address: "123 db street",
				}

				var columns = []string{`id`, `email`, `name`, `address`}
				var bln, err = Insert(db, m, columns, On{
					On: "CONFLICT(id) DO UPDATE SET email='email@db.com', name='model', address='123 db street'",
				})

				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				bln, err = Insert(db, m, columns, On{
					On: "CONFLICT(id) DO UPDATE SET email='email@db.com', name='model', address='123 db street'",
				})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				deInitDB(db)
			}
		})
	})
}
