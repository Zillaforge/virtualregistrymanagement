package server

import (
	ctlExport "VirtualRegistryManagement/controllers/grpc/export"
	ctlMemberAcl "VirtualRegistryManagement/controllers/grpc/memberacl"
	ctlOpenstack "VirtualRegistryManagement/controllers/grpc/openstack"
	ctlProject "VirtualRegistryManagement/controllers/grpc/project"
	ctlProjectAcl "VirtualRegistryManagement/controllers/grpc/projectacl"
	ctlRegistry "VirtualRegistryManagement/controllers/grpc/registry"
	ctlRepository "VirtualRegistryManagement/controllers/grpc/repository"
	ctlTag "VirtualRegistryManagement/controllers/grpc/tag"

	"google.golang.org/grpc"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

func grpcRouters(srv *grpc.Server) {
	pb.RegisterProjectCRUDControllerServer(srv, new(ctlProject.Method))
	pb.RegisterRepositoryCRUDControllerServer(srv, new(ctlRepository.Method))
	pb.RegisterTagCRUDControllerServer(srv, new(ctlTag.Method))
	pb.RegisterMemberAclCRUDControllerServer(srv, new(ctlMemberAcl.Method))
	pb.RegisterProjectAclCRUDControllerServer(srv, new(ctlProjectAcl.Method))
	pb.RegisterRegistryCRUDControllerServer(srv, new(ctlRegistry.Method))
	pb.RegisterExportCRUDControllerServer(srv, new(ctlExport.Method))
	pb.RegisterOpenstackCRUDControllerServer(srv, new(ctlOpenstack.Method))
}
