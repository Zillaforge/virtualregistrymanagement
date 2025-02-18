package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) UpdateTag(ctx context.Context, input *pb.UpdateTagInput) (output *pb.TagDetail, err error) {
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

	updateInput := &storCom.UpdateTagInput{
		ID: input.ID,
		UpdateData: &storCom.TagUpdateInfo{
			Name:            input.Name,
			ReferenceTarget: input.ReferenceTarget,
			Size:            input.Size,
			Status:          input.Status,
			Extra: func() *[]byte {
				if len(input.Extra) != 0 {
					return &input.Extra
				}
				return nil
			}(),
		},
	}

	updateOutput, err := storages.Use().UpdateTag(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageTagNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCTagNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateTag()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &updateOutput.Tag)
	return
}
