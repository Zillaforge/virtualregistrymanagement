package cmd

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/server"
	"fmt"

	"github.com/spf13/cobra"
)

func NewServe() (cmd *cobra.Command) {
	description := "Start %s Service"
	cmd = &cobra.Command{
		Use:   "serve",
		Short: fmt.Sprintf(description, cnt.UpperAbbrName),
		Long:  fmt.Sprintf(description, cnt.PascalCaseName),
		Run: func(cmd *cobra.Command, args []string) {
			server.Run()
		},
	}
	return
}
