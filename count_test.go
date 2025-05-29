package dblite

import (
	"testing"
	"time"

	"github.com/franela/goblin"
)

func TestCount(t *testing.T) {
	var g = goblin.Goblin(t)

	g.Describe("Test Count", func() {
		g.It("count: sqlite3, postgres", func() {
			g.Timeout(1 * time.Hour)
			for _, driver := range testDrivers {
				var db = initDB(driver)

				var models = []*TestModel{
					{Id: 1, Email: "email1@db.com", Name: "model1", Address: "123 db street"},
					{Id: 2, Email: "email2@db.com", Name: "model2", Address: "124 db street"},
					{Id: 3, Email: "email3@db.com", Name: "model1", Address: "125 db street"},
					{Id: 4, Email: "email4@db.com", Name: "model4", Address: "126 db street"},
					{Id: 5, Email: "email5@db.com", Name: "model1", Address: "127 db street"},
				}

				for _, model := range models {
					var bln, err = Insert(db, model, []string{
						`id`, `email`, `name`, `address`,
					}, On{On: "CONFLICT(id) DO NOTHING"})
					g.Assert(bln).IsTrue()
					g.Assert(err).IsNil()
				}

				num, err := Count(db, NewModel(-1), `id`, WhereClause{
					Where: db.SetClause("name"), Arguments: []any{"model1"},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(3))

				num, err = Count(db, NewModel(-1), `id`, WhereClause{
					Where: db.SetClause("name"), Arguments: []any{"model4"},
				})
				g.Assert(err).IsNil()
				g.Assert(num).Equal(int64(1))

				deInitDB(db)
			}
		})
	})
}
