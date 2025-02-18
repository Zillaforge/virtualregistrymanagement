package execution

import (
	"VirtualRegistryManagement/modules/openstack/cinder"
	"VirtualRegistryManagement/modules/openstack/glance"

	"github.com/gophercloud/gophercloud"
)

type AdminResource interface {
	Glance() *glance.Glance
	Cinder() *cinder.Cinder
}

type Admin struct {
	providerClient func(*string) (*gophercloud.ProviderClient, error)
	serviceClient  func(*gophercloud.ProviderClient, string) (*gophercloud.ServiceClient, error)

	glance func() *glance.Glance
	cinder func() *cinder.Cinder
}

// Glance ...
func (p *Admin) Glance() *glance.Glance {
	return p.glance()
}

// Cinder ...
func (p *Admin) Cinder() *cinder.Cinder {
	return p.cinder()
}

func (cfg *Connection) Admin() AdminResource {
	return &Admin{
		providerClient: cfg.providerClient,
		serviceClient:  cfg.serviceClient,

		glance: func() *glance.Glance {
			return cfg.Glance(cfg.AdminProject)
		},
		cinder: func() *cinder.Cinder {
			return cfg.Cinder(cfg.AdminProject)
		},
	}
}
