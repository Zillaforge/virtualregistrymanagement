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

// DeleteRepository ...
func (o *Operation) DeleteRepository(ctx context.Context, input *common.DeleteRepositoryInput) (output *common.DeleteRepositoryOutput, err error) {
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
	if deleteErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.Repository{}), &input.Where).Pluck("id", &id).Delete(&tables.Repository{}).Error; deleteErr != nil {
		if sqlErr, ok := deleteErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// Repository already by reference
			case mariadb_com.ER_ROW_IS_REFERENCED_2:
				zap.L().With(
					zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
					zap.String(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.Repository{})),
					zap.Any("value", input),
				).Warn(deleteErr.Error())
				err = tkErr.New(cnt.StorageRepositoryInUseErr).WithInner(deleteErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Repository{})),
			zap.Any("value", input),
		).Error(deleteErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(deleteErr)
		return
	}
	output = &common.DeleteRepositoryOutput{
		ID: id,
	}
	return
}
