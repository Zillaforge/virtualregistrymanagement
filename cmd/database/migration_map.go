package database

import (
	stor "VirtualRegistryManagement/storages"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Zillaforge/toolkits/mviper"
)

// MigrationMapCMD ...
func MigrationMapCMD() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "migrationmap",
		Short: "list the Database migration map",
		Long:  "list the Database migration map",
		Run: func(cmd *cobra.Command, args []string) {
			stor.New(mviper.GetString("storage.provider"))
			current, versions, err := stor.Exec().ListMigrationMap()
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(versions) == 0 {
				fmt.Println("There is no migration version in database")
				return
			}
			for _, version := range versions {
				if version == current {
					version = strings.Join([]string{version, "(It's used in currently)"}, " ")
				}
				fmt.Println(version)
			}
		},
	}
	return
}
