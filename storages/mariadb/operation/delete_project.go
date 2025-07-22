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
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

// DeleteProject ...
func (o *Operation) DeleteProject(ctx context.Context, input *common.DeleteProjectInput) (output *common.DeleteProjectOutput, err error) {
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

	var id []string
	if deleteErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.Project{}), &input.Where).Pluck("id", &id).Delete(&tables.Project{}).Error; deleteErr != nil {
		if sqlErr, ok := deleteErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// Project already by reference
			case mariadb_com.ER_ROW_IS_REFERENCED_2:
				zap.L().With(
					zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
					zap.String(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.Project{})),
					zap.Any("value", input),
				).Warn(deleteErr.Error())
				err = tkErr.New(cnt.StorageProjectInUseErr).WithInner(deleteErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Project{})),
			zap.Any("value", input),
		).Error(deleteErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(deleteErr)
		return
	}
	output = &common.DeleteProjectOutput{
		ID: id,
	}
	return
}
