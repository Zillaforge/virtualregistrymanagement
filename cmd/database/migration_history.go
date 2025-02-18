package database

import (
	stor "VirtualRegistryManagement/storages"
	"fmt"

	"github.com/spf13/cobra"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

// MigrationHistoryCMD ...
func MigrationHistoryCMD() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "migrationhistory",
		Short: "list the Database migration history",
		Long:  "list the Database migration history",
		Run: func(cmd *cobra.Command, args []string) {
			// new connection
			stor.New(mviper.GetString("storage.provider"))
			// list the migration history from database migration table
			versions, err := stor.Exec().ListMigrationHistory()
			if err != nil {
				fmt.Println(err)
			}
			if len(versions) == 0 {
				fmt.Println("There is no migration version in database")
				return
			}
			for _, version := range versions {
				fmt.Println(*version)
			}
		},
	}
	return
}
