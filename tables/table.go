package tables

import (
	"context"
	"github.com/uptrace/bun"
)

type Data struct {
	bun.BaseModel `bun:"table:data"`
	Id            int    `bun:"id,pk,autoincrement"`
	Time          string `bun:"time,notnull"`
	Type          string `bun:"type,notnull"`
	Data          string `bun:"data,notnull"`
}

func Initial(db *bun.DB, ctx context.Context) {
	_, err := db.NewCreateTable().
		Model((*Data)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		panic(err)
	}

}
