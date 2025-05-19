package dblite

import (
	"github.com/franela/goblin"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	var g = goblin.Goblin(t)
	g.Describe("Test Delete", func() {

		g.It("delete: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var data = generateData(1024)
				var bln, err = InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
					On{On: "CONFLICT(id) DO NOTHING"})
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				num, err := Count(db, NewModel(-1), `id`, WhereClause{
					Where: db.SetClause("active"), Arguments: []any{1},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(512))

				num, err = Delete(db, NewModel(-1), WhereClause{
					Where: db.SetClause("active"), Arguments: []any{1},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(512))

				num, err = Count(db, NewModel(-1), `id`, WhereClause{
					Where: db.SetClause("active"), Arguments: []any{1},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(0))

				num, err = Count(db, NewModel(-1), `id`, WhereClause{
					Where: db.SetClause("active"), Arguments: []any{0},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(512))

				deInitDB(db)
			}

		})
	})
}
