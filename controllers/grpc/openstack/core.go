package openstack

import (
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.OpenstackCRUDControllerServer
type Method struct {
	// Embed UnsafeOpenstackCRUDControllerServer to have mustEmbedUnimplementedOpenstackCRUDControllerServer()
	pb.UnsafeOpenstackCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.OpenstackCRUDControllerServer = (*Method)(nil)
