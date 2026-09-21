package dblite

import (
	"testing"
	"time"

	"github.com/franela/goblin"
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
					WithOn(
						NewOn("CONFLICT(id) DO NOTHING"),
					))
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				num, err := Count(db, NewModel(-1), `id`, Where{
					clause: db.WhereParam("active", "="), arguments: []any{1},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(512))

				num, err = Delete(db, NewModel(-1), Where{
					clause: db.WhereParam("active", "="), arguments: []any{1},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(512))

				num, err = Count(db, NewModel(-1), `id`, Where{
					clause: db.WhereParam("active", "="), arguments: []any{1},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(0))

				num, err = Count(db, NewModel(-1), `id`, Where{
					clause: db.WhereParam("active", "="), arguments: []any{0},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(512))

				deInitDB(db)
			}

		})
	})
}
