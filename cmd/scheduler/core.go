package scheduler

import (
	"VirtualRegistryManagement/cmd/args"
	cnt "VirtualRegistryManagement/constants"
	"path"

	"github.com/spf13/cobra"
)

var _schedulerConfig string = path.Join(cnt.GlobalConfigPath, "vrm-scheduler.yaml")

// NewSchedulerCmd ...
func NewSchedulerCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "scheduler",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(StartCmd())
	cmd.PersistentFlags().StringVarP(&args.CfgFileScheduler, "scheduler-config", "s", _schedulerConfig, "scheduler config file")
	return
}
