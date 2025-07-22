package memberacl

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

// Method is implement all methods as pb.MemberAclCRUDControllerServer
type Method struct {
	// Embed UnsafeMemberAclCRUDControllerServer to have mustEmbedUnimplementedMemberAclCRUDControllerServer()
	pb.UnsafeMemberAclCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.MemberAclCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(ctx context.Context, input *tables.MemberAcl) (output *pb.MemberAclDetail) {
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

	output = &pb.MemberAclDetail{
		ID: input.ID,
		Tag: &pb.TagInfo{
			ID:   input.TagID,
			Name: input.Tag.Name,
			Type: input.Tag.Type,
		},
		UserID: input.UserID,
	}
	return
}
