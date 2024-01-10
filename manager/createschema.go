package manager

import (
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// CreateSchema creates the tables for given models
func CreateSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		opt := &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}
		err := db.Model(model).CreateTable(opt)
		if err != nil {
			log.Fatal(err)
		}
	}
}
