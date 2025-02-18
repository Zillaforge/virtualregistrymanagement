package scheduler

import (
	gArgs "VirtualRegistryManagement/cmd/args"
	gComm "VirtualRegistryManagement/cmd/common"
	"VirtualRegistryManagement/configs"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/server"
	"fmt"

	"github.com/spf13/cobra"
)

func StartCmd() (cmd *cobra.Command) {
	description := "Start %s Scheduler Service"
	cmd = &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf(description, cnt.UpperAbbrName),
		Long:  fmt.Sprintf(description, cnt.PascalCaseName),
		Run: func(cmd *cobra.Command, args []string) {
			server.RunScheduler()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			configs.InitScheduler()
			gComm.MergeConfig(gArgs.CfgFileScheduler)
		},
	}
	return
}
