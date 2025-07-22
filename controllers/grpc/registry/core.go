package registry

import (
	"VirtualRegistryManagement/storages/tables"
	"time"

	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.RegistryCRUDControllerServer
type Method struct {
	// Embed UnsafeRegistryCRUDControllerServer to have mustEmbedUnimplementedRegistryCRUDControllerServer()
	pb.UnsafeRegistryCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.RegistryCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.Registry) (output *pb.RegistryInfo) {
	output = &pb.RegistryInfo{
		RepositoryID:        input.RepositoryID,
		TagID:               input.TagID,
		Creator:             input.Creator,
		ProjectID:           input.ProjectID,
		Namespace:           input.Namespace,
		RepositoryName:      input.RepositoryName,
		TagName:             input.TagName,
		Description:         input.Description,
		OperatingSystem:     input.OperatingSystem,
		Type:                input.Type,
		Size:                input.Size,
		Status:              input.Status,
		Extra:               input.Extra,
		ReferenceTarget:     input.ReferenceTarget,
		MemberAclID:         input.MemberAclID,
		AllowUserID:         input.AllowUserID,
		ProjectAclID:        input.ProjectAclID,
		AllowProjectID:      input.AllowProjectID,
		RepositoryCreatedAt: input.RepositoryCreatedAt.UTC().Format(time.RFC3339),
		RepositoryUpdatedAt: input.RepositoryUpdatedAt.UTC().Format(time.RFC3339),
		TagCreatedAt:        input.TagCreatedAt.UTC().Format(time.RFC3339),
		TagUpdatedAt:        input.TagUpdatedAt.UTC().Format(time.RFC3339),
	}
	return
}
