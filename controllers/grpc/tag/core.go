package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.TagCRUDControllerServer
type Method struct {
	// Embed UnsafeTagCRUDControllerServer to have mustEmbedUnimplementedTagCRUDControllerServer()
	pb.UnsafeTagCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.TagCRUDControllerServer = (*Method)(nil)

var empty = &emptypb.Empty{}

func (m Method) storage2pb(ctx context.Context, input *tables.Tag) (output *pb.TagDetail) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
		err       error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	if input.Repository.ID == "" {
		getInput := &storCom.GetRepositoryInput{
			ID: input.RepositoryID,
		}
		getOutput, err := storages.Use().GetRepository(ctx, getInput)
		if err != nil {
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cnt.StorageRepositoryNotFoundErr.Code():
					err = tkErr.New(cCnt.GRPCRepositoryNotFoundErr)
					return
				}
			}
			zap.L().With(
				zap.String(cnt.GRPC, "storages.Use().GetRepository()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
		input.Repository = getOutput.Repository
	}

	switch input.Status {
	case "queued", "saving", "importing", // image status
		"creating", "restoring": // snapshot status
		syncOutput, syncErr := grpc.SyncOpstkResource(ctx, &grpc.SyncOpstkResourceInput{
			Namespace:       input.Repository.Namespace,
			ProjectID:       input.Repository.ProjectID,
			TagID:           input.ID,
			Type:            input.Type,
			ReferenceTarget: input.ReferenceTarget,
		})
		if syncErr == nil {
			input = &syncOutput.Tag
		}
	case "active", "killed", "pending_delete", "deactivated": // image status
	case "available", "backing-up", "deleting", "error", "unmanaging", "error_deleting": // snapshot status
	case "deleted": // image and snapshot status
	}

	output = &pb.TagDetail{
		Tag: &pb.TagInfo{
			ID:              input.ID,
			Name:            input.Name,
			RepositoryID:    input.RepositoryID,
			ReferenceTarget: input.ReferenceTarget,
			Type:            input.Type,
			Size:            input.Size,
			Status:          input.Status,
			Extra:           input.Extra,
			CreatedAt:       input.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:       input.UpdatedAt.UTC().Format(time.RFC3339),
			Protect:         input.Protect,
		},
		Repository: &pb.RepositoryInfo{
			ID:              input.Repository.ID,
			Name:            input.Repository.Name,
			Namespace:       input.Repository.Namespace,
			OperatingSystem: input.Repository.OperatingSystem,
			Description:     input.Repository.Description,
			Creator:         input.Repository.Creator,
			ProjectID:       input.Repository.ProjectID,
			CreatedAt:       input.Repository.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:       input.Repository.UpdatedAt.UTC().Format(time.RFC3339),
			Protect:         input.Repository.Protect,
		},
	}

	return
}
