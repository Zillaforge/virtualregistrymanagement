package database

import "github.com/spf13/cobra"

// NewDatabaseCmd ...
func NewDatabaseCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "database",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(SyncCmd(), MigrationMapCMD(), MigrationHistoryCMD(), MigrateCmd(), RollbackCmd())
	return
}
