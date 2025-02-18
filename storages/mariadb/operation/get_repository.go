package operation

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// GetRepository ...
func (o *Operation) GetRepository(ctx context.Context, input *common.GetRepositoryInput) (output *common.GetRepositoryOutput, err error) {
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

	output = &common.GetRepositoryOutput{}
	if getErr := o.conn.WithContext(ctx).Preload("Tag").Model(&tables.Repository{}).Where("id = ?", input.ID).First(&output.Repository).Error; getErr != nil {
		if errors.Is(getErr, gorm.ErrRecordNotFound) {
			// Repository not found
			zap.L().With(
				zap.String(cnt.Storage, "record not found"),
				zap.String(cnt.RequestID, requestID),
				zap.String("name", GetTableName(&tables.Repository{})),
			).Error(getErr.Error())
			err = tkErr.New(cnt.StorageRepositoryNotFoundErr).WithInner(getErr)
			return
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).First(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Repository{})),
		).Error(getErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(getErr)
		return
	}
	return
}
