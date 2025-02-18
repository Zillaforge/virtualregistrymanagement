package operation

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages/common"
	mariadb_com "VirtualRegistryManagement/storages/mariadb/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// DeleteProjectAcl ...
func (o *Operation) DeleteProjectAcl(ctx context.Context, input *common.DeleteProjectAclInput) (output *common.DeleteProjectAclOutput, err error) {
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

	var tagID []string
	if deleteErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.ProjectAcl{}), &input.Where).Pluck("tag_id", &tagID).Delete(&tables.ProjectAcl{}).Error; deleteErr != nil {
		if sqlErr, ok := deleteErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// ProjectAcl already by reference
			case mariadb_com.ER_ROW_IS_REFERENCED_2:
				zap.L().With(
					zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
					zap.String(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.ProjectAcl{})),
					zap.Any("value", input),
				).Warn(deleteErr.Error())
				err = tkErr.New(cnt.StorageProjectAclInUseErr).WithInner(deleteErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.ProjectAcl{})),
			zap.Any("value", input),
		).Error(deleteErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(deleteErr)
		return
	}
	output = &common.DeleteProjectAclOutput{
		TagID: tagID,
	}
	return
}
