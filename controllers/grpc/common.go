package grpc

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack"
	"VirtualRegistryManagement/modules/openstack/common"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"VirtualRegistryManagement/utility/querydecoder"
	"VirtualRegistryManagement/utility/workerpool"
	"context"
	"reflect"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
)

// A generic empty message that you can re-use to avoid defining duplicated
var EmptyPb = &emptypb.Empty{}

func WhereErrorParser(input error) error {
	return whereErrorParser(reflect.ValueOf(input))
}
func whereErrorParser(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			return whereErrorParser(reflect.ValueOf(v.MapIndex(key).Interface()))
		}
	case reflect.Struct:
		switch err := v.Interface().(type) {
		case querydecoder.UnknownKeyError:
			return tkErr.New(cCnt.GRPCWhereBindingErr).With("field", err.Key)
		case querydecoder.RegexError:
			return tkErr.New(cCnt.GRPCWhereBindingErr).WithInner(err)
		}
	}
	return nil
}

type (
	SyncOpstkResourceInput struct {
		Namespace       string
		ProjectID       string
		TagID           string
		Type            string
		ReferenceTarget string
	}

	OpstkResourceInput struct {
		Registries *storCom.ListRegistriesOutput

		Creator      *string
		ProjectID    *string
		RepositoryID *string
		TagID        *string
	}

	OpstkImageInput struct {
		Namespace string
		ProjectID string
		ImageID   string

		RepositoryID *string
		TagID        *string
	}

	OpstkSnapshotInput struct {
		Namespace  string
		ProjectID  string
		SnapshotID string

		RepositoryID *string
		TagID        *string
	}
)

const (
	TagTypeImage          = "common"
	TagTypeVolumeSnapshot = "increase"
)

func SyncOpstkResource(ctx context.Context, input *SyncOpstkResourceInput) (output *storCom.UpdateTagOutput, err error) {
	var (
		requestID = utility.MustGetContextRequestID(ctx)
		status    string
		size      uint64
	)

	switch input.Type {
	case TagTypeImage:
		getInput := &common.GetImageInput{ID: input.ReferenceTarget}
		getOutput, _err := openstack.Namespace(input.Namespace).Glance(input.ProjectID).GetImage(ctx, getInput)
		if _err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Glance().GetImage()"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", getOutput),
			).Error(_err.Error())
			err = _err
			return
		}
		status = getOutput.Image.Status.String()
		size = uint64(getOutput.Image.SizeBytes)

	case TagTypeVolumeSnapshot:
		getInput := &common.GetSnapshotInput{ID: input.ReferenceTarget}
		getOutput, _err := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).GetSnapshot(ctx, getInput)
		if _err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Cinder().GetSnapshot()"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", getOutput),
			).Error(_err.Error())
			err = _err
			return
		}
		status = getOutput.Snapshot.Status.String()
		size = uint64(getOutput.Snapshot.SizeBytes)
	}

	updateInput := &storCom.UpdateTagInput{
		ID: input.TagID,
		UpdateData: &storCom.TagUpdateInfo{
			Status: &status,
			Size:   &size,
		},
	}
	updateOutput, _err := storages.Use().UpdateTag(ctx, updateInput)
	if _err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(_err.Error())
		err = _err
		return
	}
	output = updateOutput
	return
}

func DeleteOpstkResource(ctx context.Context, input *OpstkResourceInput) {
	requestID := utility.MustGetContextRequestID(ctx)

	var (
		namespace string
		projectID string
	)

	type (
		tags struct {
			tag   map[string]bool
			typ   string
			total int64

			repositoryID *string
			tagID        *string
		}
	)

	target := map[string]*tags{}
	for _, registry := range input.Registries.Registries {
		namespace = registry.Namespace
		projectID = registry.ProjectID
		resourceID := registry.ReferenceTarget

		if _, ok := target[resourceID]; !ok {
			listTagsInput := &storCom.ListTagsInput{
				Pagination: storCom.Paginate(-1, 0),
				Where: storCom.ListTagWhere{
					ReferenceTarget: &resourceID,
				},
			}
			listTagsOutput, err := storages.Use().ListTags(ctx, listTagsInput)
			if err != nil {
				zap.L().With(
					zap.String(cnt.GRPC, "storages.Use().ListTags(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", listTagsInput),
				).Error(err.Error())
				continue
			}

			target[resourceID] = &tags{
				tag:   map[string]bool{},
				typ:   registry.Type,
				total: listTagsOutput.Count,

				repositoryID: input.RepositoryID,
				tagID:        &registry.TagID,
			}
			for _, tag := range listTagsOutput.Tags {
				target[resourceID].tag[tag.ID] = false
			}
		}

		target[resourceID].tag[registry.TagID] = true
	}

	for resourceID, tag := range target {
		isAll := true
		for _, t := range tag.tag {
			isAll = isAll && t
		}

		switch isAll {
		case false: // only remove related information in resource tags/metadata
			switch tag.typ {
			case TagTypeImage:
				updateImageInput := OpstkImageInput{
					Namespace: namespace,
					ProjectID: projectID,
					ImageID:   resourceID,

					RepositoryID: tag.repositoryID,
					TagID:        tag.tagID,
				}
				for !workerpool.Use().TrySubmit(UpdateOpstkImage(requestID, updateImageInput)) {
				}
			case TagTypeVolumeSnapshot:
				updateSnapshotInput := OpstkSnapshotInput{
					Namespace:  namespace,
					ProjectID:  projectID,
					SnapshotID: resourceID,

					RepositoryID: tag.repositoryID,
					TagID:        tag.tagID,
				}
				for !workerpool.Use().TrySubmit(UpdateOpstkSnapshot(requestID, updateSnapshotInput)) {
				}
			}

		case true: // delete resource in openstack
			switch tag.typ {
			case TagTypeImage:
				deleteImageInput := OpstkImageInput{
					Namespace: namespace,
					ProjectID: projectID,
					ImageID:   resourceID,
				}
				for !workerpool.Use().TrySubmit(DeleteOpstkImage(requestID, deleteImageInput)) {
				}
			case TagTypeVolumeSnapshot:
				deleteSnapshotInput := OpstkSnapshotInput{
					Namespace:  namespace,
					ProjectID:  projectID,
					SnapshotID: resourceID,
				}
				for !workerpool.Use().TrySubmit(DeleteOpstkSnapshot(requestID, deleteSnapshotInput)) {
				}
			}
		}
	}
}

func UpdateOpstkImage(requestID string, input OpstkImageInput) func() {
	ctx := tracer.StartEntryContext(requestID)

	zap.L().With(
		zap.String(cnt.GRPC, "UpdateOpstkImage()"),
		zap.String(cnt.RequestID, requestID),
		zap.String("namespace", input.Namespace),
		zap.String("project-id", input.ProjectID),
	).Info(input.ImageID)

	return func() {
		getImageInput := &opstkCom.GetImageInput{
			ID: input.ImageID,
		}
		getImageOutput, err := openstack.Namespace(input.Namespace).Glance(input.ProjectID).GetImage(ctx, getImageInput)
		if err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Glance().GetImage(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", getImageInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
		updateImageInput := &opstkCom.UpdateImageInput{
			ID:                input.ImageID,
			Tags:              &getImageOutput.Image.Tags,
			DeleteSystemLabel: true,

			RepositoryID: input.RepositoryID,
			TagID:        input.TagID,
		}
		if _, err := openstack.Namespace(input.Namespace).Glance(input.ProjectID).UpdateImage(ctx, updateImageInput); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Glance().UpdateImage(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", updateImageInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
	}
}

func DeleteOpstkImage(requestID string, input OpstkImageInput) func() {
	ctx := tracer.StartEntryContext(requestID)

	zap.L().With(
		zap.String(cnt.GRPC, "DeleteOpstkImage()"),
		zap.String(cnt.RequestID, requestID),
		zap.String("namespace", input.Namespace),
		zap.String("project-id", input.ProjectID),
	).Info(input.ImageID)

	return func() {
		deleteImageInput := &opstkCom.DeleteImageInput{
			ID: input.ImageID,
		}
		if _, err := openstack.Namespace(input.Namespace).Admin().Glance().DeleteImage(ctx, deleteImageInput); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Admin().Glance().DeleteImage(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", deleteImageInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
	}
}

func UpdateOpstkSnapshot(requestID string, input OpstkSnapshotInput) func() {
	ctx := tracer.StartEntryContext(requestID)

	zap.L().With(
		zap.String(cnt.GRPC, "UpdateOpstkSnapshot()"),
		zap.String(cnt.RequestID, requestID),
		zap.String("namespace", input.Namespace),
		zap.String("project-id", input.ProjectID),
	).Info(input.SnapshotID)

	return func() {
		getSnapshotInput := &opstkCom.GetSnapshotInput{
			ID: input.SnapshotID,
		}
		getSnapshotOutput, err := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).GetSnapshot(ctx, getSnapshotInput)
		if err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Glance().GetSnapshot(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", getSnapshotInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
		updateSnapshotMetadataInput := &opstkCom.UpdateSnapshotMetadataInput{
			ID:                input.SnapshotID,
			Metadata:          &getSnapshotOutput.Snapshot.Metadata,
			DeleteSystemLabel: true,

			RepositoryID: input.RepositoryID,
			TagID:        input.TagID,
		}
		if _, err := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).UpdateSnapshotMetadata(ctx, updateSnapshotMetadataInput); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Cinder().DeleteSnapshot(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", updateSnapshotMetadataInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
	}
}

func DeleteOpstkSnapshot(requestID string, input OpstkSnapshotInput) func() {
	ctx := tracer.StartEntryContext(requestID)

	zap.L().With(
		zap.String(cnt.GRPC, "DeleteOpstkSnapshot()"),
		zap.String(cnt.RequestID, requestID),
		zap.String("namespace", input.Namespace),
		zap.String("project-id", input.ProjectID),
	).Info(input.SnapshotID)

	return func() {
		deleteSnapshotInput := &opstkCom.DeleteSnapshotInput{
			ID: input.SnapshotID,
		}
		if _, err := openstack.Namespace(input.Namespace).Admin().Cinder().DeleteSnapshot(ctx, deleteSnapshotInput); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Admin().Cinder().DeleteSnapshot(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", deleteSnapshotInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
	}
}
