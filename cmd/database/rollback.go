package database

import (
	stor "VirtualRegistryManagement/storages"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Zillaforge/toolkits/mviper"
)

const syncMessage string = "Syncing %s database"

// RollbackCmd ...
func RollbackCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "rollback",
		Short: "rollback the Database migration by ID",
		Long:  "rollback the Database to specified migration ID",
		Run: func(cmd *cobra.Command, args []string) {
			stor.New(mviper.GetString("storage.provider"))
			if err := stor.Exec().Rollback(); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Database rollback success")

		},
	}
	return
}
