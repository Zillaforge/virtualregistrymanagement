package cmd

import (
	"VirtualRegistryManagement/cmd/args"
	"VirtualRegistryManagement/cmd/database"
	"VirtualRegistryManagement/cmd/scheduler"
	cnt "VirtualRegistryManagement/constants"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/mviper"
)

var (
	rootCmd = &cobra.Command{
		Use:   cnt.Kind,
		Short: cnt.UpperAbbrName,
		Long:  cnt.Name,
	}
	globalConfig = path.Join(cnt.GlobalConfigPath, cnt.GlobalConfigFilename)
)

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(NewVersion(), NewServe(), database.NewDatabaseCmd(), scheduler.NewSchedulerCmd())
	rootCmd.PersistentFlags().StringVarP(&args.CfgFileG, "config", "c", globalConfig, "config file")
}

func initConfig() {
	mviper.SetConfigType("yaml")
	if args.CfgFileG != "" {
		mviper.SetConfigFile(args.CfgFileG)
	}

	if err := mviper.MergeInConfig(); err != nil {
		err = tkErr.New(cnt.ConfigReadConfigFailedErr).WithInner(err)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := mviper.VerifyConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if mviper.GetString("version") != cnt.Version {
		err := tkErr.New(cnt.ConfigVersionMustBeSameErr)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
