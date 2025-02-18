package cinder

import (
	cnt "VirtualRegistryManagement/constants"

	"github.com/gophercloud/gophercloud"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

type Cinder struct {
	namespace string
	projectID string
	sc        *gophercloud.ServiceClient
}

func New(namespace, projectID string, sc *gophercloud.ServiceClient) *Cinder {
	return &Cinder{
		namespace: namespace,
		projectID: projectID,
		sc:        sc,
	}
}

func (c *Cinder) SetServiceClient(namespace string, sc *gophercloud.ServiceClient) *Cinder {
	c.namespace = namespace
	c.sc = sc
	return c
}

func (c *Cinder) checkConnection() error {
	if c.sc == nil {
		return tkErr.New(cnt.OpenstackConnectionIsNotCreatedErr)
	}
	return nil
}
