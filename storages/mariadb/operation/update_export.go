package operation

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// UpdateExport ...
func (o *Operation) UpdateExport(ctx context.Context, input *common.UpdateExportInput) (output *common.UpdateExportOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"error":  &err,
	})

	whereCondition := &common.UpdateExportInput{
		ID: input.ID,
	}
	output = &common.UpdateExportOutput{}
	if updateErr := o.conn.WithContext(ctx).Model(&tables.Export{}).Where(queryConversion(*whereCondition)).Updates(queryConversion(*input.UpdateData)).First(&output.Export).Error; err != nil {
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).Updates(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Export{})),
			zap.Any("value", input),
		).Error(updateErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(updateErr)
		return
	}
	return
}
