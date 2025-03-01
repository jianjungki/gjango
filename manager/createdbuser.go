package manager

import (
	"fmt"

	"tiktok_tools/config"

	"github.com/go-pg/pg/v10"
)

// CreateDatabaseUserIfNotExist creates a database user
func CreateDatabaseUserIfNotExist(db *pg.DB, p *config.PostgresConfig) {
	statement := fmt.Sprintf(`SELECT * FROM pg_roles WHERE rolname = '%s';`, p.User)
	res, err := db.Exec(statement)
	if err != nil {
		fmt.Println("creating database error" + err.Error())
	}
	if res.RowsReturned() == 0 {
		statement = fmt.Sprintf(`CREATE USER %s WITH PASSWORD '%s';`, p.User, p.Password)
		_, err := db.Exec(statement)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf(`Created user %s`, p.User)
		}
	} else {
		fmt.Printf("Database user %s already exists\n", p.User)
	}

}
