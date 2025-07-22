package database

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Zillaforge/toolkits/mviper"

	stor "VirtualRegistryManagement/storages"
)

var (
	id string
)

// MigrateCmd ...
func MigrateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate the Database version",
		Long:  "migrate the Database to specified version",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				version := args[0]
				stor.New(mviper.GetString("storage.provider"))
				if err := stor.Exec().Migrate(version); err != nil {
					panic(err)
				}
				fmt.Printf("Database is migrated to version %s\n", version)
			} else {
				fmt.Println("You need to provide the version number which you want to migrate to ...")
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&id, "id", "i", "0.2.0", "migration id (default is 0.2.0")
	cmd.MarkPersistentFlagRequired("version")
	return
}
