package openstack

import (
	"VirtualRegistryManagement/authentication"
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/cinder"
	opstkExec "VirtualRegistryManagement/modules/openstack/execution"
	"VirtualRegistryManagement/modules/openstack/glance"
	"VirtualRegistryManagement/services"
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/flatten"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

const _openstackKind string = "openstack"

var provider map[string]Resource = map[string]Resource{}

type Resource interface {
	Admin() opstkExec.AdminResource
	Glance(projectID string) *glance.Glance
	Cinder(projectID string) *cinder.Cinder
}

func Init(config interface{}) (err error) {
	for _, service := range cast.ToSlice(config) {
		if service == nil {
			continue
		}
		cfg := viper.New()
		cfg.MergeConfigMap(cast.ToStringMap(service))

		if !cfg.IsSet("namespace") || !cfg.IsSet("service") {
			return tkErr.New(cnt.OpenstackNamespaceAndServiceIsRequiredErr)
		}

		namespace := cfg.GetString("namespace")
		service := cfg.GetString("service")

		if _, exist := provider[namespace]; exist {
			return tkErr.New(cnt.OpenstackNamespaceIsRepeatedErr)
		}

		switch services.ServiceMap[service].Kind {
		case _openstackKind:
			if value, ok := services.ServiceMap[service].Conn.(*opstkExec.Connection); ok {
				p := opstkExec.Connection{
					IdentityEndpoint: value.IdentityEndpoint,
					Username:         value.Username,
					Password:         value.Password,
					DomainName:       value.DomainName,
					AllowReauth:      value.AllowReauth,
					AdminProject:     value.AdminProject,
					PidSource:        value.PidSource,
					Pid:              GetOpstkPID(value.PidSource),
				}
				p.SetNamespace(namespace)
				provider[namespace] = &p
			}
		default:
			return tkErr.New(cnt.OpenstackTypeIsNotSupportedErr)
		}
	}
	return nil
}

func Namespace(namespace string) Resource {
	return provider[namespace]
}

func NamespaceIsLegal(namespace string) bool {
	_, ok := provider[namespace]
	return ok
}

func ListNamespaces() []string {
	namespaces := []string{}
	for ns := range provider {
		namespaces = append(namespaces, ns)
	}
	return namespaces
}

// Replace replaces global provider by p
func Replace(namespace string, p opstkExec.Connection) {
	provider[namespace] = &p
}

func GetOpstkPID(pidSource string) func(string) string {
	ctx := tracer.StartEntryContext(tracer.EmptyRequestID)

	return func(projectID string) string {
		authProjectInput := &authCom.GetProjectInput{ID: projectID, Cacheable: true}
		authProjectOutput, _ := authentication.Use().GetProject(ctx, authProjectInput)
		projectInfo, _ := flatten.Flatten(authProjectOutput.ToMap(), "", flatten.DotStyle)
		return fmt.Sprintf("%v", projectInfo[pidSource])
	}
}
