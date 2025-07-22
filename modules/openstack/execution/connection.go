package execution

import (
	cnt "VirtualRegistryManagement/constants"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	tkErr "github.com/Zillaforge/toolkits/errors"
)

type Connection struct {
	IdentityEndpoint string
	Username         string
	Password         string
	DomainName       string
	AllowReauth      bool
	AdminProject     string
	PidSource        string
	Pid              func(string) string

	namespace string
}

func New(p *Connection) (*Connection, error) {
	_, err := p.providerClient(nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (cfg *Connection) providerClient(projectID *string) (*gophercloud.ProviderClient, error) {
	authOptions := gophercloud.AuthOptions{
		IdentityEndpoint: cfg.IdentityEndpoint,
		Username:         cfg.Username,
		Password:         cfg.Password,
		DomainName:       cfg.DomainName,
		AllowReauth:      cfg.AllowReauth,
	}
	if projectID != nil {
		authOptions.Scope = &gophercloud.AuthScope{
			ProjectID: *projectID,
		}
	}
	return openstack.AuthenticatedClient(authOptions)
}

func (Connection) serviceClient(pc *gophercloud.ProviderClient, key string) (*gophercloud.ServiceClient, error) {
	var scFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	switch key {
	case _glanceResource:
		scFunc = openstack.NewImageServiceV2
	case _cinderResource:
		scFunc = openstack.NewBlockStorageV3
	default:
		return nil, tkErr.New(cnt.OpenstackTypeIsNotSupportedErr)
	}
	return scFunc(pc, gophercloud.EndpointOpts{})
}

func (cfg *Connection) SetNamespace(namespace string) {
	cfg.namespace = namespace
}
