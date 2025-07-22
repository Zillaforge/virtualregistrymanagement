package eventpublish

import (
	cnt "VirtualRegistryManagement/constants"
	epCom "VirtualRegistryManagement/eventpublish/common"
	"VirtualRegistryManagement/eventpublish/grpc"
	"VirtualRegistryManagement/eventpublish/native"
	"VirtualRegistryManagement/utility"
	"strings"

	"gopkg.in/yaml.v3"

	"fmt"
	"os"
	"reflect"

	EB "github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/mviper"
)

var (
	actionBus  EB.Bus = EB.New()
	asyncBus   EB.Bus = EB.New()
	syncBus    EB.Bus = EB.New()
	plugins           = map[string]epCom.EventIntf{}
	priorities        = map[string]int{}
)

func reconcile(action string, meta map[string]string, req interface{}, resp interface{}) {
	asyncBus.Publish(cnt.AsyncKey, action, meta, req, resp)
	syncBus.Publish(cnt.SyncKey, action, meta, req, resp)
}

func InitAllPlugins() {
	for _, plugin := range mviper.Get("event_publish.plugins").([]interface{}) {
		configMap := interface2map(plugin)
		configYaml, _ := yaml.Marshal(plugin)
		if t, ok := configMap["type"].(string); ok {
			name, ok := configMap["name"].(string)
			if !ok {
				continue
			}
			switch t {
			case epCom.GRPCType:
				startPluginInput := &grpc.StartPluginInput{
					Name: name,
					Type: epCom.GRPCType,
				}
				for _, key := range []string{
					"binary_path",
					"config_path",
					"socket_path",
				} {
					path, ok := configMap[key].(string)
					if !ok {
						// TODO
						return
					}
					switch key {
					case "binary_path":
						startPluginInput.BinaryPath = path
					case "config_path":
						startPluginInput.ConfigPath = path
					case "socket_path":
						startPluginInput.SocketPath = path
					}
				}

				if args, ok := configMap["start_cmd_args"].(map[string]interface{}); ok {
					startPluginInput.StartCmdArgs = args
				}

				// 啟動 Plugin Binary
				core, err := grpc.StartPlugin(startPluginInput)
				if err != nil {
					zap.L().With(
						zap.String(cnt.Plugin, "grpc.StartPlugin(...)"),
						zap.Any("input", startPluginInput),
					).Error(err.Error())
					os.Exit(1)
				}

				plugins[name] = core
			case epCom.NativeType:
				if plugin := native.EnableNativePlugins(name); plugin != nil {
					plugins[name] = plugin
				}
			default:
				continue
			}

			// 設定 Native Sample Plugin Configuration
			if plugins[name] != nil {
				plugins[name].SetConfig(configYaml)
			}

			// 設定非同步優先順序
			if priority := configMap["priority"]; priority != nil {
				if value, ok := priority.(int); ok {
					priorities[name] = value
				}
			}
		}
	}
	logstr := []string{}
	for key, plugin := range plugins {
		// 初始化 Plugin
		if !plugin.InitPlugin() {
			zap.L().With(
				zap.String(cnt.Plugin, "if !plugin.InitPlugin() ..."),
			).Error(cnt.PluginInternalServerErrMsg)
			os.Exit(2)
		}

		// 確認配置與服務版本是否一致
		if !plugin.CheckPluginVersion() {
			zap.L().With(
				zap.String(cnt.Plugin, "if !plugin.CheckPluginVersion() ..."),
			).Error(fmt.Sprintf(cnt.PluginVersionDoesNotMatchErrMsg, key))
			os.Exit(1)
		}

		// 啟動非同步訂閱
		if _, ok := priorities[key]; !ok {
			zap.L().Info(fmt.Sprint("SubscribeAsync ", key))
			asyncBus.SubscribeAsync(cnt.AsyncKey, plugin.Reconcile, false)
		}
		logstr = append(logstr, key)
	}

	// 根據 priority 調整同步訂閱順序
	pl := utility.SortValues(priorities)

	// 啟動同步訂閱
	for _, pair := range pl {
		zap.L().Info(fmt.Sprint("Subscribe ", pair.Key))
		syncBus.Subscribe(cnt.SyncKey, plugins[pair.Key].Reconcile)
	}

	// subscribe reconcile function
	actionBus.SubscribeAsync(cnt.ReconcileKey, reconcile, false)
	zap.L().Info(fmt.Sprint("Loaded plugins (", len(plugins), "): ", strings.Join(logstr, ", ")))
}

// GetBus ...
func GetBus() EB.Bus {
	return actionBus
}

func interface2map(input interface{}) (m map[string]interface{}) {
	if m == nil {
		m = make(map[string]interface{})
	}

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			if k, ok := e.Interface().(string); ok {
				m[k] = val.MapIndex(e).Interface()
			}
		}
	}
	return
}

func EnableHTTPRouters(rg *gin.RouterGroup) {
	for _, plugin := range plugins {
		plugin.EnableHTTPRouter(rg.Group(strings.ToLower(plugin.GetName())))
	}
}
