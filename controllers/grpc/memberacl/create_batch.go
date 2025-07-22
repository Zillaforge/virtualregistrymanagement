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

func (m *Method) CreateMemberAclBatch(ctx context.Context, input *pb.MemberAclBatchInfo) (output *pb.MemberAclBatchDetail, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	output = &pb.MemberAclBatchDetail{
		Data: []*pb.MemberAclDetail{},
	}

	listInput := &pb.ListNamespaceInput{}
	listOutput, err := m.ListMemberAcls(ctx, listInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "m.ListMemberAcl"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	existAcl := map[string]bool{}
	for _, memberAcl := range listOutput.Data {
		if memberAcl.Tag != nil {
			existAcl[memberAcl.Tag.ID+memberAcl.UserID] = true
		}
	}

	createInput := &storCom.CreateMemberAclBatchInput{}
	for _, data := range input.Data {
		if existAcl[data.TagID+data.UserID] {
			continue
		}
		createInput.MemberAcls = append(createInput.MemberAcls, tables.MemberAcl{
			TagID:  data.TagID,
			UserID: data.UserID,
		})
	}

	if len(createInput.MemberAcls) == 0 {
		return
	}

	createOutput, err := storages.Use().CreateMemberAclBatch(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageMemberAclNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCMemberAclNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateMemberAclBatch()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	for _, memberAcl := range createOutput.MemberAcls {
		output.Data = append(output.Data, m.storage2pb(ctx, &memberAcl))
	}
	return
}
