package export

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

func (m *Method) GetExport(ctx context.Context, input *pb.GetInput) (output *pb.ExportInfo, err error) {
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

	getInput := &storCom.GetExportInput{
		ID: input.ID,
	}

	getOutput, err := storages.Use().GetExport(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageExportNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCExportNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetExport()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &getOutput.Export)
	return
}
