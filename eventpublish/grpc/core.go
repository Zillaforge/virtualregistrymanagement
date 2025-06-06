package grpc

import (
	cnt "VirtualRegistryManagement/constants"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

type StartPluginInput struct {
	Name         string
	Type         string
	BinaryPath   string
	ConfigPath   string
	SocketPath   string
	StartCmdArgs map[string]interface{}
	_            struct{}
}

var _plugins []*exec.Cmd

/*
StartPlugin ...

errors:
- 99000000(internal server error)
*/
func StartPlugin(input *StartPluginInput) (c *core, err error) {
	c = &core{
		name: input.Name,
		t:    input.Type,
	}
	binaryExist, configExist := true, true
	if _, err := os.Stat(input.BinaryPath); errors.Is(err, os.ErrNotExist) {
		zap.L().Debug(fmt.Sprintf("Binary file not found: '%s'", input.BinaryPath))
		binaryExist = false
	}
	if _, err := os.Stat(input.ConfigPath); errors.Is(err, os.ErrNotExist) {
		zap.L().Debug(fmt.Sprintf("Config file not found: '%s'", input.ConfigPath))
		configExist = false
	}
	os.Remove(input.SocketPath)
	// 需要同時滿足 binary file, config file 皆存在且 sock file 不存在的條件下
	// 系統才會自動以子程序方式啟動 plugin process。
	if binaryExist && configExist {
		zap.L().Info(fmt.Sprintf("Plugin %s is starting ...", input.Name))
		cmd := exec.Command(input.BinaryPath, cmdArgs(input.ConfigPath, input.StartCmdArgs, "serve")...)
		if err := cmd.Start(); err != nil {
			zap.L().Error(err.Error())
			return nil, tkErr.New(cnt.PluginInternalServerErr)
		}
		_plugins = append(_plugins, cmd)
		c.cmd = cmd

		defer func() {
			go startMonitor(c)
		}()
	}

	newGRPCConnInput := &newGRPCConnInput{
		conn:       5,
		socketPath: input.SocketPath,
	}
	c.handler = newGRPCConn(newGRPCConnInput)
	return c, nil
}

func cmdArgs(configPath string, argMap map[string]interface{}, subCmd string) (args []string) {
	args = append(args, "-c", configPath)

	for key, value := range argMap {
		args = append(args, fmt.Sprintf("%v", key), fmt.Sprintf("%v", value))
	}

	args = append(args, subCmd)
	return
}

func startMonitor(c *core) {
	for {
		if c != nil && c.GetName() == "" {
			zap.L().Warn(fmt.Sprintf("%s plugin has broken ...", c.name))
			if err := c.cmd.Start(); err != nil {
				c.cmd.Process = nil
				zap.L().Warn(fmt.Sprintf("Create %s plugin again. But it's failed ...", c.name))
			} else {
				zap.L().Info(fmt.Sprintf("Create %s plugin is successful", c.name))
			}
		}
		time.Sleep(120 * time.Second)
	}
}

func ClosePlugins() {
	for _, plugin := range _plugins {
		plugin.Wait()
	}
}
