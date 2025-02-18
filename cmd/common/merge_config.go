package common

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"
	"os"

	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

func MergeConfig(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Can't find Configuration (%s): %v\n", path, err)
		os.Exit(1)
	}

	v := mviper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Parser Global Configuration Failed : %v at: %s \n", err, path)
		os.Exit(1)
	}

	if err := mviper.MergeConfigMap(v.AllSettings()); err != nil {
		fmt.Printf("Parser Global Configuration Failed : %v at: %s \n", err, path)
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
