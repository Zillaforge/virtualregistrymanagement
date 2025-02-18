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

// ListProjectAcls ...
func (o *Operation) ListProjectAcls(ctx context.Context, input *common.ListProjectAclsInput) (output *common.ListProjectAclsOutput, err error) {
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

	output = &common.ListProjectAclsOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}
	if listErr := whereCascade(o.conn.WithContext(ctx).Preload("Tag").Preload("ProjectAcl").Model(&tables.Registry{}), &input.Where).Where("project_acl_id is not null").
		Where("project_acl_id is not null").Count(&output.Count).Limit(limit).Offset(offset).Find(&output.ProjectAcls).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Preload(...).Preload(...).Model(...), ...).Count(...).Limit(...).Offset(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.ProjectAcl{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
