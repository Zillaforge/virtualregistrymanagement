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

// DeleteMemberAcl ...
func (o *Operation) DeleteMemberAcl(ctx context.Context, input *common.DeleteMemberAclInput) (output *common.DeleteMemberAclOutput, err error) {
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
	if deleteErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.MemberAcl{}), &input.Where).Pluck("id", &id).Delete(&tables.MemberAcl{}).Error; deleteErr != nil {
		if sqlErr, ok := deleteErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// MemberAcl already by reference
			case mariadb_com.ER_ROW_IS_REFERENCED_2:
				zap.L().With(
					zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
					zap.String(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.MemberAcl{})),
					zap.Any("value", input),
				).Warn(deleteErr.Error())
				err = tkErr.New(cnt.StorageMemberAclInUseErr).WithInner(deleteErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.MemberAcl{})),
			zap.Any("value", input),
		).Error(deleteErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(deleteErr)
		return
	}
	output = &common.DeleteMemberAclOutput{
		ID: id,
	}
	return
}
