package dblite

import (
	"testing"
	"time"

	"github.com/franela/goblin"
)

func TestUpdate(t *testing.T) {
	var g = goblin.Goblin(t)
	g.Describe("Test Update", func() {

		g.It("update: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var data = generateData(100)
				var bln, err = InsertMany(db, data, []string{`id`, `email`, `name`, `address`, `active`},
					WithOn(
						NewOn("CONFLICT(id) DO NOTHING"),
					))
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				var email = "db-email-73@db.com"
				var address = "73 street, ghana"
				var model = data[73]
				model.Email = email
				model.Address = address
				model.Active = 0

				cols, err := ColumnsByExclusion(NewModel(-1), []string{"id", "active"})
				g.Assert(err).IsNil()

				bln, err = Update(db, model, cols, WithWhere(NewWhere(
					db.SetParam("id", len(cols)+1), WithWhereArguments([]any{model.Id}),
				)))
				g.Assert(bln).IsTrue()
				g.Assert(err).IsNil()

				object, err := QueryModel(db, model, WithWhere(NewWhere(
					db.SetParam("id"), WithWhereArguments([]any{model.Id}),
				)))
				g.Assert(object.Id).Equal(model.Id)
				g.Assert(object.Email).Equal(email)
				g.Assert(object.Address).Equal(address)
				g.Assert(err).IsNil()

				deInitDB(db)
			}

		})
	})
}
