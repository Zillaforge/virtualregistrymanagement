package system

import (
	"VirtualRegistryManagement/utility"
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"github.com/Zillaforge/toolkits/configs"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtil "github.com/Zillaforge/toolkits/utilities"
)

// GetSystemConfigurationsOutput ...
type GetSystemConfigurationsOutput struct {
	Configurations []Configuration `json:"configurations"`
	UnknownKeys    []string        `json:"unknownKeys"`
}

// Configuration ...
type Configuration struct {
	Key          string      `json:"key"`
	Value        interface{} `json:"value"`
	DefaultValue interface{} `json:"defaultValue"`
	TypeOf       string      `json:"typeOf"`
}

// GetSystemConfigurations ...
func GetSystemConfigurations(c *gin.Context) {
	var (
		funcName   = tkUtil.NameOfFunction().String()
		statusCode = http.StatusOK
		output     = &GetSystemConfigurationsOutput{}
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"output":     output,
		"statusCode": &statusCode,
	})

	output = &GetSystemConfigurationsOutput{Configurations: make([]Configuration, 0), UnknownKeys: make([]string, 0)}
	k := mviper.AllKeys()
	sort.Strings(k)
	for _, key := range k {
		cfg := configs.Get(key)
		if cfg == nil {
			output.UnknownKeys = append(output.UnknownKeys, key)
			continue
		}
		switch t := mviper.Get(cfg.Key).(type) {
		case []interface{}:
			if len(t) > 0 {
				switch t[0].(type) {
				case int:
					output.Configurations = append(output.Configurations, Configuration{
						Key:          cfg.Key,
						Value:        fmt.Sprintf("%+v", t),
						TypeOf:       "[]int",
						DefaultValue: fmt.Sprintf("%+v", cfg.DefaultValue),
					})
				case string:
					output.Configurations = append(output.Configurations, Configuration{
						Key:          cfg.Key,
						Value:        fmt.Sprintf("%+v", t),
						TypeOf:       "[]string",
						DefaultValue: fmt.Sprintf("%+v", cfg.DefaultValue),
					})
				case map[interface{}]interface{}:
					for idx, itm := range t {
						y, _ := yaml.Marshal(itm)
						output.Configurations = append(output.Configurations, Configuration{
							Key:          fmt.Sprintf("%s.%d", cfg.Key, idx),
							Value:        string(y),
							TypeOf:       "yaml",
							DefaultValue: "none",
						})
					}

				default:
					output.Configurations = append(output.Configurations, Configuration{
						Key:          cfg.Key,
						Value:        t,
						TypeOf:       "unknown",
						DefaultValue: cfg.DefaultValue,
					})
				}
			} else {
				output.Configurations = append(output.Configurations, Configuration{
					Key:          cfg.Key,
					Value:        fmt.Sprintf("%+v", t),
					TypeOf:       "[]interface{}",
					DefaultValue: cfg.DefaultValue,
				})
			}
		default:
			output.Configurations = append(output.Configurations, Configuration{
				Key:          cfg.Key,
				Value:        t,
				TypeOf:       reflect.TypeOf(mviper.Get(cfg.Key)).Name(),
				DefaultValue: cfg.DefaultValue,
			})
		}

	}
	utility.ResponseWithType(c, http.StatusOK, output)
}
