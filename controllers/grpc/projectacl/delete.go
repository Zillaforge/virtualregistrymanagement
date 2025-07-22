package projectacl

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/modules/openstack"
	opstkComm "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"VirtualRegistryManagement/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

func (m *Method) DeleteProjectAcl(ctx context.Context, input *pb.DeleteInput) (output *pb.DeleteOutput, err error) {
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

	// binding the where parameter
	whereInput := storCom.DeleteProjectAclWhere{}
	if err = querydecoder.ShouldBindWhereSlice(&whereInput, input.Where); err != nil {
		if e, ok := tkErr.IsError(grpc.WhereErrorParser(err)); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				return output, e
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "querydecoder.ShouldBindWhereSlice(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input.Where", input.Where),
		).Error(err.Error())
		return output, tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
	}

	deleteInput := &storCom.DeleteProjectAclInput{
		Where: whereInput,
	}

	deleteOutput, err := storages.Use().DeleteProjectAcl(ctx, deleteInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectAclNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCProjectAclNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteProjectAcl()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	updateOpstkVisibilityToPrivate := func(namespace, projectID, imageID string) (err error) {
		publicVisibility := opstkComm.ImageVisibilityPrivate
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

	for _, tagID := range deleteOutput.TagID {
		listProjectAclsInput := &pb.ListNamespaceInput{
			Limit:  -1,
			Offset: 0,
			Where:  []string{"tag-id=" + tagID},
		}
		listProjectAclsOutput, listProjectAclsErr := m.ListProjectAcls(ctx, listProjectAclsInput)
		if listProjectAclsErr != nil {

		}
		if listProjectAclsOutput.Count != 0 { // still share
			continue
		}

		getTagInput := &pb.GetInput{
			ID: tagID,
		}
		getTagOutput, getTagErr := vrm.GetTag(getTagInput, ctx)
		if getTagErr != nil {
			// Expected errors
			if e, ok := tkErr.IsError(getTagErr); ok {
				switch e.Code() {
				case cCnt.GRPCTagNotFoundErr.Code():
					err = e
					return
				}
			}
			// Unexpected errors
			zap.L().With(
				zap.String(cnt.GRPC, "vrm.GetTag()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getTagInput),
			).Error(getTagErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getTagErr)
			return
		}
		if getTagOutput.Tag.Type == grpc.TagTypeImage { // update image visibility
			var (
				namespace = getTagOutput.Repository.Namespace
				projectID = getTagOutput.Repository.ProjectID
				imageID   = getTagOutput.Tag.ReferenceTarget
			)
			if err = updateOpstkVisibilityToPrivate(namespace, projectID, imageID); err != nil {
				return
			}
		}
	}

	output = &pb.DeleteOutput{
		ID: deleteOutput.TagID,
	}
	return
}
