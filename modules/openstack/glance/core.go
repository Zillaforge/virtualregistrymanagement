package glance

import (
	cnt "VirtualRegistryManagement/constants"

	"github.com/gophercloud/gophercloud"
	tkErr "github.com/Zillaforge/toolkits/errors"
)

type Glance struct {
	namespace string
	projectID string
	sc        *gophercloud.ServiceClient
}

func New(namespace, projectID string, sc *gophercloud.ServiceClient) *Glance {
	return &Glance{
		namespace: namespace,
		projectID: projectID,
		sc:        sc,
	}
}

func (g *Glance) SetServiceClient(namespace string, sc *gophercloud.ServiceClient) *Glance {
	g.namespace = namespace
	g.sc = sc
	return g
}

func (g *Glance) checkConnection() error {
	if g.sc == nil {
		return tkErr.New(cnt.OpenstackConnectionIsNotCreatedErr)
	}
	return nil
}
