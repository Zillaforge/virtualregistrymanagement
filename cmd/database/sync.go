package database

import (
	stor "VirtualRegistryManagement/storages"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

// SyncCmd ...
func SyncCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "sync",
		Short: "sync the Database",
		Long:  "create the Database table",
		Run: func(cmd *cobra.Command, args []string) {
			provider := mviper.GetString("storage.provider")
			dbname := mviper.GetString(fmt.Sprintf("storage.%s.name", provider))
			zap.L().Info(fmt.Sprintf(syncMessage, dbname))
			stor.New(provider)
			fmt.Println("Sync database is successful !!")
		},
	}
	return
}
