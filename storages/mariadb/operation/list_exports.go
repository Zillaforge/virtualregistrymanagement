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

// ListExports ...
func (o *Operation) ListExports(ctx context.Context, input *common.ListExportsInput) (output *common.ListExportsOutput, err error) {
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

	output = &common.ListExportsOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}
	if listErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.Export{}), &input.Where).Count(&output.Count).Limit(limit).Offset(offset).Find(&output.Exports).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Count(...).Limit(...).Offset(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Export{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
