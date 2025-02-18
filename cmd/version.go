package cmd

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// NewVersion ...
func NewVersion() (cmd *cobra.Command) {
	description := "Show %s Version"
	cmd = &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf(description, cnt.UpperAbbrName),
		Long:  fmt.Sprintf(description, cnt.PascalCaseName),
		Run: func(cmd *cobra.Command, args []string) {
			basicVersion()
		},
	}
	return
}

func basicVersion() {
	fmt.Printf("%s: %v\n", cnt.Kind, cnt.Version)
	fmt.Printf("Golang: %v\n", runtime.Version())
}
