package dblite

import (
	"testing"
	"time"

	"github.com/franela/goblin"
)

func TestQuery(t *testing.T) {
	var g = goblin.Goblin(t)

	g.Describe("Test Query", func() {

		g.It("model: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var data = generateData(100)
				var bln, err = InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				model, err := QueryModel(db, NewModel(-1), WhereClause{
					Where: db.SetParam("id"), Arguments: []any{73},
				})
				g.Assert(model.Id).Eql(int64(73))
				g.Assert(err).IsNil()

				deInitDB(db)
			}

		})

		g.It("models: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var data = generateData(1024)
				var bln, err = InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				results, err := QueryModels(db, NewModel(-1), WhereClause{Where: `"active"=1`})
				g.Assert(len(results)).Eql(512)
				g.Assert(err).IsNil()

				results, err = QueryModels(db, NewModel(-1), WhereClause{Where: `"active"=0`})
				g.Assert(len(results)).Eql(512)
				g.Assert(err).IsNil()

				results, err = QueryModels(db, NewModel(-1), WhereClause{
					Where: db.WhereParam("active", "="), Arguments: []any{2},
				})
				g.Assert(len(results)).Eql(0)
				g.Assert(err).IsNil()

				results, err = QueryModels(db, NewModel(-1), WhereClause{
					Where: db.WhereParam("id", "="), Arguments: []any{1096},
				})
				g.Assert(len(results)).Eql(0)
				g.Assert(err).IsNil()

				deInitDB(db)
			}

		})
	})
}
