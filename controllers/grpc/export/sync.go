package export

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/filepath"
	"VirtualRegistryManagement/modules/openstack"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"VirtualRegistryManagement/utility/workerpool"
	"context"

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

func (m *Method) SyncExports(ctx context.Context, input *emptypb.Empty) (output *emptypb.Empty, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)
	output = empty

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"output": &output,
			"error":  &err,
		},
	)

	listInput := &pb.ListNamespaceInput{
		Limit:  -1,
		Offset: 0,
		Where: []string{
			"status=preparing",
		},
	}
	listOutput, err := m.ListExports(ctx, listInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.EventConsume, "vrm.ListExports(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(err.Error())
		return
	}

	for _, export := range listOutput.Data {
		var imageStatus *string
		if export.Status != opstkCom.ImageStatusActive.String() {
			getImageInput := &opstkCom.GetImageInput{
				ID: export.ImageID,
			}
			getImageOutput, getImageErr := openstack.Namespace(export.Namespace).Glance(export.ProjectID).GetImage(ctx, getImageInput)
			if getImageErr != nil {
				zap.L().With(
					zap.String(cnt.Controller, "openstack.Namespace().Glance().GetImage(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", getImageInput),
				).Error(getImageErr.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getImageErr)
				continue
			}

			switch status := getImageOutput.Image.Status.String(); status {
			case opstkCom.ImageStatusActive.String():
				imageStatus = &status
			case export.Status:
				continue
			default:
				updateExportInput := &pb.UpdateExportInput{
					ID:          export.ID,
					ImageStatus: &status,
				}
				if _, updateExportErr := vrm.UpdateExport(updateExportInput, ctx); err != nil {
					zap.L().With(
						zap.String(cnt.Controller, "vrm.UpdateExport(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.String("namespace", export.Namespace),
						zap.String("project-id", export.ProjectID),
						zap.Any("input", updateExportInput),
					).Error(updateExportErr.Error())
					err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(updateExportErr)
				}
				continue
			}
		}

		status := "exporting"
		updateExportInput := &pb.UpdateExportInput{
			ID:          export.ID,
			ImageStatus: imageStatus,
			Status:      &status,
		}

		if _, updateExportErr := vrm.UpdateExport(updateExportInput, ctx); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "vrm.UpdateExport(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", export.Namespace),
				zap.String("project-id", export.ProjectID),
				zap.Any("input", updateExportInput),
			).Error(updateExportErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(updateExportErr)
			continue
		}

		path, err := filepath.Path(export.Filepath, []interface{}{export.ProjectID}, nil)
		if err != nil {
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cnt.FilepathTypeIsNotSupportedErr.Code():
					err = tkErr.New(cCnt.GRPCNotSupportTypeErr)
					continue
				}
			}
			zap.L().With(
				zap.String(cnt.Controller, "filepath.Path()"),
				zap.String(cnt.RequestID, requestID),
				zap.String("project-id", export.ProjectID),
				zap.String("filepath", export.Filepath),
			).Info(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			continue
		}

		downloadImageInput := DownloadImageInput{
			Namespace: export.Namespace,
			ProjectID: export.ProjectID,
			TagID:     export.TagID,
			ImageID:   export.ImageID,
			Filepath:  path,
			ExportID:  export.ID,
		}
		for !workerpool.Use().TrySubmit(downloadImage(requestID, downloadImageInput)) {
		}
	}

	return
}

func downloadImage(requestID string, input DownloadImageInput) func() {
	ctx := tracer.StartEntryContext(requestID)

	zap.L().With(
		zap.String(cnt.GRPC, "downloadImage()"),
		zap.String(cnt.RequestID, requestID),
		zap.String("namespace", input.Namespace),
		zap.String("project-id", input.ProjectID),
		zap.String("tag-id", input.TagID),
		zap.String("image-id", input.ImageID),
	).Info(input.Filepath)

	return func() {
		downloadImageInput := &opstkCom.DownloadImageDataInput{
			ID:       input.ImageID,
			Filepath: input.Filepath,
		}
		if _, err := openstack.Namespace(input.Namespace).Glance(input.ProjectID).DownloadImageData(ctx, downloadImageInput); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Glance().DownloadImageData(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", downloadImageInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}

		status := "finished"
		updateExportInput := &pb.UpdateExportInput{
			ID:     input.ExportID,
			Status: &status,
		}
		updateExportOutput, err := vrm.UpdateExport(updateExportInput, ctx)
		if err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "vrm.UpdateExport(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", updateExportInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}

		switch updateExportOutput.Type {
		case "image":
			// do nothing
		case "snapshot":
			deleteVolumeInput := &opstkCom.DeleteVolumeInput{
				ID: *updateExportOutput.VolumeID,
			}
			_, deleteVolumeErr := openstack.Namespace(updateExportOutput.Namespace).Cinder(updateExportOutput.ProjectID).DeleteVolume(ctx, deleteVolumeInput)
			if deleteVolumeErr != nil {
				zap.L().With(
					zap.String(cnt.Controller, "openstack.Namespace().Cinder().DeleteVolume(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.String("namespace", input.Namespace),
					zap.String("project-id", input.ProjectID),
					zap.Any("input", deleteVolumeInput),
				).Error(deleteVolumeErr.Error())
			}

			deleteImageInput := &opstkCom.DeleteImageInput{
				ID: updateExportOutput.ImageID,
			}
			_, deleteImageErr := openstack.Namespace(updateExportOutput.Namespace).Glance(updateExportOutput.ProjectID).DeleteImage(ctx, deleteImageInput)
			if deleteImageErr != nil {
				zap.L().With(
					zap.String(cnt.Controller, "openstack.Namespace().Glance().DeleteImage(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.String("namespace", input.Namespace),
					zap.String("project-id", input.ProjectID),
					zap.Any("input", deleteImageInput),
				).Error(deleteImageErr.Error())
			}
		}
	}
}
