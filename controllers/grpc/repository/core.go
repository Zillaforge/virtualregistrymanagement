package repository

import (
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/storages/tables"
	"context"
	"time"

	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.RepositoryCRUDControllerServer
type Method struct {
	// Embed UnsafeRepositoryCRUDControllerServer to have mustEmbedUnimplementedRepositoryCRUDControllerServer()
	pb.UnsafeRepositoryCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.RepositoryCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(ctx context.Context, input *tables.Repository) (output *pb.RepositoryDetail) {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	output = &pb.RepositoryDetail{
		Repository: &pb.RepositoryInfo{
			ID:              input.ID,
			Name:            input.Name,
			Namespace:       input.Namespace,
			OperatingSystem: input.OperatingSystem,
			Description:     input.Description,
			Creator:         input.Creator,
			ProjectID:       input.ProjectID,
			CreatedAt:       input.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:       input.UpdatedAt.UTC().Format(time.RFC3339),
			Protect:         input.Protect,
		},
		Tags: []*pb.TagInfo{},
	}

	for _, tag := range input.Tag {
		switch tag.Status {
		case "queued", "saving", "importing", // image status
			"creating", "restoring": // snapshot status
			syncOutput, syncErr := grpc.SyncOpstkResource(ctx, &grpc.SyncOpstkResourceInput{
				Namespace:       input.Namespace,
				ProjectID:       input.ProjectID,
				TagID:           tag.ID,
				Type:            tag.Type,
				ReferenceTarget: tag.ReferenceTarget,
			})
			if syncErr == nil {
				tag = syncOutput.Tag
			}
		case "active", "killed", "pending_delete", "deactivated": // image status
		case "available", "backing-up", "deleting", "error", "unmanaging", "error_deleting": // snapshot status
		case "deleted": // image and snapshot status
		}

		output.Tags = append(output.Tags, &pb.TagInfo{
			ID:              tag.ID,
			Name:            tag.Name,
			RepositoryID:    tag.RepositoryID,
			ReferenceTarget: tag.ReferenceTarget,
			Type:            tag.Type,
			Size:            tag.Size,
			Status:          tag.Status,
			Extra:           tag.Extra,
			CreatedAt:       tag.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:       tag.UpdatedAt.UTC().Format(time.RFC3339),
			Protect:         input.Protect,
		})
	}

	return
}
