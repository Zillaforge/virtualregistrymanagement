package projectacl

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/modules/openstack"
	opstkComm "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

func (m *Method) CreateProjectAclBatch(ctx context.Context, input *pb.ProjectAclBatchInfo) (output *pb.ProjectAclBatchDetail, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	output = &pb.ProjectAclBatchDetail{
		Data: []*pb.ProjectAclDetail{},
	}

	updateOpstkVisibilityToPublic := func(namespace, projectID, imageID string) (err error) {
		publicVisibility := opstkComm.ImageVisibilityPublic
		updateInput := &opstkComm.UpdateImageInput{
			ID:         imageID,
			Visibility: &publicVisibility,
		}
		_, err = openstack.Namespace(namespace).Glance(projectID).UpdateImage(ctx, updateInput)
		if err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Glance().UpdateImage(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", namespace),
				zap.Any("input", updateInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
		return
	}

	createInput := &storCom.CreateProjectAclBatchInput{}
	for _, data := range input.Data {
		getTagInput := &pb.GetInput{
			ID: data.TagID,
		}
		getTagOutput, getTagErr := vrm.GetTag(getTagInput, ctx)
		if getTagErr != nil {
			// Expected errors
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cCnt.GRPCTagNotFoundErr.Code():
					err = e
					return
				}
			}
			// Unexpected errors
			zap.L().With(
				zap.String(cnt.GRPC, "vrm.GetInput(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getTagInput),
			).Error(getTagErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getTagErr)
			return
		}

		var count int64 = 0
		if data.ProjectID == nil { // global public
			// delete global limit acl, because public
			deleteProjectAclInput := &pb.DeleteInput{
				Where: []string{"tag-id=" + getTagOutput.Tag.ID},
			}
			_, deleteProjectAclErr := m.DeleteProjectAcl(ctx, deleteProjectAclInput)
			if deleteProjectAclErr != nil {
				zap.L().With(
					zap.String(cnt.GRPC, "vrm.DeleteProjectAcl(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", deleteProjectAclInput),
				).Error(deleteProjectAclErr.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(deleteProjectAclErr)
				return
			}

		} else { // global limit
			listProjectAclsInput := &pb.ListNamespaceInput{
				Limit:  -1,
				Offset: 0,
				Where:  []string{"tag-id=" + getTagOutput.Tag.ID},
			}
			listProjectAclsOutput, listProjectAclsErr := m.ListProjectAcls(ctx, listProjectAclsInput)
			if listProjectAclsErr != nil {
				zap.L().With(
					zap.String(cnt.GRPC, "vrm.ListProjectAcls(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", listProjectAclsInput),
				).Error(listProjectAclsErr.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(listProjectAclsErr)
				return
			}

			for _, projectAcl := range listProjectAclsOutput.Data {
				if projectAcl.ProjectID == nil { // already public
					continue
				}
			}
			count = listProjectAclsOutput.Count
		}

		if count == 0 &&
			getTagOutput.Tag.Type == grpc.TagTypeImage { // update image visibility
			var (
				namespace = getTagOutput.Repository.Namespace
				projectID = getTagOutput.Repository.ProjectID
				imageID   = getTagOutput.Tag.ReferenceTarget
			)
			if err = updateOpstkVisibilityToPublic(namespace, projectID, imageID); err != nil {
				return
			}
		}

		createInput.ProjectAcls = append(createInput.ProjectAcls, tables.ProjectAcl{
			TagID:     data.TagID,
			ProjectID: data.ProjectID,
		})
	}

	if len(createInput.ProjectAcls) == 0 {
		return
	}

	createOutput, err := storages.Use().CreateProjectAclBatch(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectAclNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCProjectAclNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateProjectAclBatch()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	for _, projectAcl := range createOutput.ProjectAcls {
		output.Data = append(output.Data, m.storage2pb(ctx, &projectAcl))
	}
	return
}
