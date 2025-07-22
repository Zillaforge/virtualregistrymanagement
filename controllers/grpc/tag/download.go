package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/modules/filepath"
	"VirtualRegistryManagement/modules/openstack"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type DownloadImageInput struct {
	Namespace string
	ProjectID string
	TagID     string
	ImageID   string
	Filepath  string
	ExportID  string
}

func (m *Method) DownloadTag(ctx context.Context, input *pb.ImageDataInput) (output *emptypb.Empty, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)
	output = empty

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	getTagInput := &pb.GetInput{
		ID: input.TagID,
	}
	getTagOutput, err := m.GetTag(ctx, getTagInput)
	if err != nil {
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
			zap.String(cnt.GRPC, "vrm.GetTag()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getTagInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	projectID := getTagOutput.Repository.ProjectID
	path, err := filepath.Validate(input.Filepath, []interface{}{projectID}, nil,
		"img", fmt.Sprintf("%s-%s", getTagOutput.Repository.Name, getTagOutput.Tag.Name))
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.FilepathTypeIsNotSupportedErr.Code():
				err = tkErr.New(cCnt.GRPCNotSupportTypeErr)
				return
			case cnt.FilepathFilepathIsExistErr.Code():
				err = tkErr.New(cCnt.GRPCFileExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "filepath.Path()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("project-id", projectID),
			zap.String("filepath", input.Filepath),
		).Info(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	namespace := getTagOutput.Repository.Namespace
	createExportInput := &pb.ExportInfo{
		Namespace:      namespace,
		ProjectID:      projectID,
		Creator:        getTagOutput.Repository.Creator,
		RepositoryID:   getTagOutput.Repository.ID,
		RepositoryName: getTagOutput.Repository.Name,
		TagID:          getTagOutput.Tag.ID,
		TagName:        getTagOutput.Tag.Name,

		Filepath: path,
		Status:   "preparing",
	}

	switch getTagOutput.Tag.Type {
	case grpc.TagTypeVolumeSnapshot:
		createExportInput.Type = "snapshot"

		getSnapshotInput := &opstkCom.GetSnapshotInput{
			ID: getTagOutput.Tag.ReferenceTarget,
		}
		getSnapshotOutput, getSnapshotErr := openstack.Namespace(namespace).Cinder(projectID).GetSnapshot(ctx, getSnapshotInput)
		if getSnapshotErr != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Cinder().GetSnapshot(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getSnapshotOutput),
			).Error(getSnapshotErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getSnapshotErr)
			return
		}
		createExportInput.SnapshotID = ToPtr(getSnapshotOutput.Snapshot.ID)
		createExportInput.SnapshotStatus = ToPtr(getSnapshotOutput.Snapshot.Status.String())

		// create volume
		createVolumeInput := &opstkCom.CreateVolumeInput{
			SnapshotID:  getSnapshotOutput.Snapshot.ID,
			Size:        int(getSnapshotOutput.Snapshot.FullSizeBytes / 1024 / 1024 / 1024),
			Name:        getSnapshotOutput.Snapshot.Name,
			Description: fmt.Sprintf("create from snapshot %s", getSnapshotOutput.Snapshot.ID),
			// Metadata:    input.Metadata,
		}
		createVolumeOutput, createVolumeErr := openstack.Namespace(namespace).Cinder(projectID).CreateVolume(ctx, createVolumeInput)
		if createVolumeErr != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Cinder().CreateVolume(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", createVolumeInput),
			).Error(createVolumeErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(createVolumeErr)
			return
		}

		volumeID := createVolumeOutput.Volume.ID
		volumeStatus := createVolumeOutput.Volume.Status
		for volumeStatus != opstkCom.VolumeStatusAvailable {
			// create volume
			getVolumeInput := &opstkCom.GetVolumeInput{
				ID: volumeID,
			}
			getOutput, getVolumeErr := openstack.Namespace(namespace).Cinder(projectID).GetVolume(ctx, getVolumeInput)
			if getVolumeErr != nil {
				zap.L().With(
					zap.String(cnt.Controller, "openstack.Namespace().Cinder().GetVolume(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", getVolumeInput),
				).Error(getVolumeErr.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getVolumeErr)
				return
			}

			volumeID = getOutput.Volume.ID
			volumeStatus = getOutput.Volume.Status
		}

		createExportInput.VolumeID = ToPtr(createVolumeOutput.Volume.ID)
		createExportInput.VolumeStatus = ToPtr(volumeStatus.String())

		// upload volume to image
		uploadImageInput := &opstkCom.UploadImageFromVolumeInput{
			ID:   createVolumeOutput.Volume.ID,
			Name: createVolumeOutput.Volume.Name,
		}
		uploadImageOutput, uploadImageErr := openstack.Namespace(namespace).Cinder(projectID).UploadImageFromVolume(ctx, uploadImageInput)
		if uploadImageErr != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Cinder().UploadImageFromVolume(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", uploadImageInput),
			).Error(uploadImageErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(uploadImageErr)
			return
		}

		createExportInput.ImageID = uploadImageOutput.Image.ImageID
		createExportInput.ImageStatus = uploadImageOutput.Image.Status.String()

	case grpc.TagTypeImage:
		createExportInput.Type = "image"

		getImageInput := &opstkCom.GetImageInput{
			ID: getTagOutput.Tag.ReferenceTarget,
		}
		getImageOutput, getImageErr := openstack.Namespace(namespace).Glance(projectID).GetImage(ctx, getImageInput)
		if getImageErr != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Glance().GetImage(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getImageInput),
			).Error(getImageErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getImageErr)
			return
		}

		createExportInput.ImageID = getImageOutput.Image.ID
		createExportInput.ImageStatus = getImageOutput.Image.Status.String()
	}

	createExportOutput, err := vrm.CreateExport(createExportInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateExport(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createExportOutput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}
	return
}

func ToPtr[T any](in T) *T {
	return &in
}
