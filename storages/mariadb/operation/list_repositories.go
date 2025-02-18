package operation

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// ListRepositories ...
func (o *Operation) ListRepositories(ctx context.Context, input *common.ListRepositoriesInput) (output *common.ListRepositoriesOutput, err error) {
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

	output = &common.ListRepositoriesOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}
	if listErr := whereCascade(o.conn.WithContext(ctx).Preload("Tag").Model(&tables.Repository{}), &input.Where).Count(&output.Count).Limit(limit).Offset(offset).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&output.Repositories).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Count(...).Limit(...).Offset(...).Order(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Repository{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
