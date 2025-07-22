package operation

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

// UpdateRepository ...
func (o *Operation) UpdateRepository(ctx context.Context, input *common.UpdateRepositoryInput) (output *common.UpdateRepositoryOutput, err error) {
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

	whereCondition := &common.UpdateRepositoryInput{
		ID: input.ID,
	}
	output = &common.UpdateRepositoryOutput{}
	if updateErr := o.conn.WithContext(ctx).Model(&tables.Repository{}).Where(queryConversion(*whereCondition)).Updates(queryConversion(*input.UpdateData)).First(&output.Repository).Error; err != nil {
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).Updates(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Repository{})),
			zap.Any("value", input),
		).Error(updateErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(updateErr)
		return
	}
	return
}
