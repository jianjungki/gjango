package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"tiktok_tools/migration"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all migrations",
	Long:  `reset all migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset called")
		err := migration.Run("reset")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	migrateCmd.AddCommand(resetCmd)
}
