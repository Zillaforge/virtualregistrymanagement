package projectacl

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.ProjectAclCRUDControllerServer
type Method struct {
	// Embed UnsafeProjectAclCRUDControllerServer to have mustEmbedUnimplementedProjectAclCRUDControllerServer()
	pb.UnsafeProjectAclCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ProjectAclCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(ctx context.Context, input *tables.ProjectAcl) (output *pb.ProjectAclDetail) {
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

	if input.Tag.ID == "" {
		getInput := &storCom.GetTagInput{
			ID: input.TagID,
		}
		getOutput, _err := storages.Use().GetTag(ctx, getInput)
		if _err != nil {
			if e, ok := tkErr.IsError(_err); ok {
				switch e.Code() {
				case cnt.StorageTagNotFoundErr.Code():
					err = tkErr.New(cCnt.GRPCTagNotFoundErr)
					return
				}
			}
			zap.L().With(
				zap.String(cnt.GRPC, "storages.Use().GetTag()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getInput),
			).Error(_err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
			return
		}
		input.Tag = getOutput.Tag
	}

	output = &pb.ProjectAclDetail{
		ID: input.ID,
		Tag: &pb.TagInfo{
			ID:   input.TagID,
			Name: input.Tag.Name,
			Type: input.Tag.Type,
		},
		ProjectID: input.ProjectID,
	}
	return
}
