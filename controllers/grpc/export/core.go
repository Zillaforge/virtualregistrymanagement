package export

import (
	"VirtualRegistryManagement/storages/tables"
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.ExportCRUDControllerServer
type Method struct {
	// Embed UnsafeExportCRUDControllerServer to have mustEmbedUnimplementedExportCRUDControllerServer()
	pb.UnsafeExportCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ExportCRUDControllerServer = (*Method)(nil)

var empty = &emptypb.Empty{}

func (m Method) storage2pb(ctx context.Context, input *tables.Export) (output *pb.ExportInfo) {
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

	output = &pb.ExportInfo{
		ID:             input.ID,
		RepositoryID:   input.RepositoryID,
		RepositoryName: input.RepositoryName,
		TagID:          input.TagID,
		TagName:        input.TagName,
		Type:           input.Type,
		SnapshotID:     input.SnapshotID,
		SnapshotStatus: input.SnapshotStatus,
		VolumeID:       input.VolumeID,
		VolumeStatus:   input.VolumeStatus,
		ImageID:        input.ImageID,
		ImageStatus:    input.ImageStatus,
		Filepath:       input.Filepath,
		Status:         input.Status,
		Creator:        input.Creator,
		ProjectID:      input.ProjectID,
		Namespace:      input.Namespace,
		CreatedAt:      input.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      input.UpdatedAt.UTC().Format(time.RFC3339),
	}
	return
}
