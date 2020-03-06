package cmd

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/manager"
	"github.com/spf13/cobra"
)

// createCmd represents the migrate command
var createdbCmd = &cobra.Command{
	Use:   "createdb",
	Short: "createdb creates a database user and database from database parameters declared in config",
	Long:  `createdb creates a database user and database from database parameters declared in config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createdb called")
		p := config.GetPostgresConfig()

		// connection to db as postgres superuser
		dbSuper := config.GetPostgresSuperUserConnection()
		defer dbSuper.Close()

		manager.CreateDatabaseUserIfNotExist(dbSuper, p)
		manager.CreateDatabaseIfNotExist(dbSuper, p)
	},
}

func init() {
	rootCmd.AddCommand(createdbCmd)
}
