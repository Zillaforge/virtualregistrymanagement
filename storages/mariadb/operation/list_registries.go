package operation

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

// ListRegistries ...
func (o *Operation) ListRegistries(ctx context.Context, input *common.ListRegistriesInput) (output *common.ListRegistriesOutput, err error) {
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

	output = &common.ListRegistriesOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}

	tx := o.conn
	// Belong the Specific User
	if input.Flag.BelongUser && input.Flag.UserID != nil {
		tx = tx.Or("creator = ?", *input.Flag.UserID)
	}
	// Belong the Specific Project
	if input.Flag.BelongProject && input.Flag.ProjectID != nil {
		tx = tx.Or("project_id = ?", *input.Flag.ProjectID)
	}
	// Share to Specific User
	if input.Flag.ProjectLimit && input.Flag.UserID != nil {
		tx = tx.Or("allow_user_id = ?", *input.Flag.UserID)
	}
	// Public in Project
	if input.Flag.ProjectPublic && input.Flag.ProjectID != nil {
		tx = tx.Or("project_id = ? AND allow_project_id = ?", *input.Flag.ProjectID, *input.Flag.ProjectID)
	}
	// Share to Specific Project
	if input.Flag.GlobalLimit && input.Flag.ProjectID != nil {
		tx = tx.Or("allow_project_id = ?", *input.Flag.ProjectID)
	}
	// Public in Namespace
	if input.Flag.GlobalPublic {
		tx = tx.Or("project_acl_id IS NOT null AND allow_project_id IS null")
	}

	if listErr := whereCascade(o.conn.WithContext(ctx).
		Preload("Project").Preload("Repository").Preload("Tag").Preload("MemberAcl").Preload("ProjectAcl").
		Model(&tables.Registry{}).Where(tx), &input.Where).Count(&output.Count).Limit(limit).Offset(offset).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "tag_created_at"}, Desc: true}).Find(&output.Registries).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...).Where(), ...).Count(...).Limit(...).Offset(...).Order(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Registry{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
