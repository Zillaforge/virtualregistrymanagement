package project

import (
	"VirtualRegistryManagement/storages/tables"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.ProjectCRUDControllerServer
type Method struct {
	// Embed UnsafeProjectCRUDControllerServer to have mustEmbedUnimplementedProjectCRUDControllerServer()
	pb.UnsafeProjectCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ProjectCRUDControllerServer = (*Method)(nil)

var empty = &emptypb.Empty{}

func (m Method) storage2pb(ctx context.Context, input *tables.Project) (output *pb.ProjectInfo) {
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

	output = &pb.ProjectInfo{
		ID:             input.ID,
		LimitCount:     input.LimitCount,
		LimitSizeBytes: input.LimitSizeBytes,
	}
	return
}
