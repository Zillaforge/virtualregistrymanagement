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
	"gorm.io/gorm/clause"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

// CreateRepository ...
func (o *Operation) CreateRepository(ctx context.Context, input *common.CreateRepositoryInput) (output *common.CreateRepositoryOutput, err error) {
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

	if createErr := o.conn.WithContext(ctx).Clauses(clause.Returning{}).Create(&input.Repository).Error; createErr != nil {
		if sqlErr, ok := createErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// Repository 已經存在
			case mariadb_com.ER_DUP_ENTRY:
				zap.L().With(
					zap.String(cnt.Storage, "o.conn.WithContext().Clauses(...).Create(...).Error"),
					zap.Any(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.Repository{})),
					zap.Any("value", input),
				).Warn(createErr.Error())
				err = tkErr.New(cnt.StorageRepositoryExistErr).WithInner(createErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Clauses(...).Create(...).Error"),
			zap.Any(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Repository{})),
			zap.Any("value", input),
		).Error(createErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(createErr)
		return
	}
	output = &common.CreateRepositoryOutput{
		Repository: input.Repository,
	}
	return
}
