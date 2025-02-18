package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/filepath"
	"VirtualRegistryManagement/modules/openstack"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"VirtualRegistryManagement/utility/workerpool"
	"context"
	"errors"
	"os"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type UploadImageInput struct {
	Namespace string
	ProjectID string
	TagID     string
	ImageID   string
	Filepath  string
}

func (m *Method) UploadTag(ctx context.Context, input *pb.ImageDataInput) (output *emptypb.Empty, err error) {
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
	path, err := filepath.Path(input.Filepath, []interface{}{projectID}, nil)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.FilepathTypeIsNotSupportedErr.Code():
				err = tkErr.New(cCnt.GRPCNotSupportTypeErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "filepath.Path()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("project-id", projectID),
			zap.String("filepath", input.Filepath),
		).Info(err.Error())
		err = tkErr.New(cCnt.GRPCNotSupportTypeErr).WithInner(err)
		return
	}

	// check file exist
	if _, err = os.Stat(path); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "os.Stat()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("filepath", path),
		).Info(err.Error())
		if errors.Is(err, os.ErrNotExist) {
			err = tkErr.New(cCnt.GRPCFileNotFoundErr)
			return
		}
		err = tkErr.New(cCnt.GRPCFileNotFoundErr).WithInner(err)
		return
	}

	uploadImageInput := UploadImageInput{
		Namespace: getTagOutput.Repository.Namespace,
		ProjectID: projectID,
		TagID:     getTagOutput.Tag.ID,
		ImageID:   getTagOutput.Tag.ReferenceTarget,
		Filepath:  path,
	}
	for !workerpool.Use().TrySubmit(uploadImage(requestID, uploadImageInput)) {
	}

	return
}

func uploadImage(requestID string, input UploadImageInput) func() {
	ctx := tracer.StartEntryContext(requestID)

	zap.L().With(
		zap.String(cnt.GRPC, "uploadImage()"),
		zap.String(cnt.RequestID, requestID),
		zap.String("namespace", input.Namespace),
		zap.String("project-id", input.ProjectID),
		zap.String("tag-id", input.TagID),
		zap.String("image-id", input.ImageID),
	).Info(input.Filepath)

	return func() {
		uploadImageInput := &common.UploadImageDataInput{
			ID:       input.ImageID,
			Filepath: input.Filepath,
		}
		if _, err := openstack.Namespace(input.Namespace).Glance(input.ProjectID).UploadImageData(ctx, uploadImageInput); err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Glance().UploadImageData(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", uploadImageInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}

		getImageInput := &common.GetImageInput{
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

		status := getImageOutput.Image.Status.String()
		size := uint64(getImageOutput.Image.SizeBytes)
		updateTagInput := &pb.UpdateTagInput{
			ID:     input.TagID,
			Size:   &size,
			Status: &status,
		}
		if _, err = vrm.UpdateTag(updateTagInput, ctx); err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "vrm.UpdateTag()"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", input.Namespace),
				zap.String("project-id", input.ProjectID),
				zap.Any("input", updateTagInput),
			).Error(err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
			return
		}
	}
}
