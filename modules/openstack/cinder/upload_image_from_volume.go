package cinder

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (c Cinder) UploadImageFromVolume(ctx context.Context, input *common.UploadImageFromVolumeInput) (output *common.UploadImageFromVolumeOutput, err error) {
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

	if err = c.checkConnection(); err != nil {
		zap.L().With(
			zap.String(cnt.Module, "g.checkConnection()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
		).Error(err.Error())
		return
	}

	uploadInput := volumeactions.UploadImageOpts{
		ImageName: input.Name,
	}
	uploadOutput, err := volumeactions.UploadImage(c.sc, input.ID, uploadInput).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "volumeactions.UploadImage(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.Any("input", uploadInput),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.UploadImageFromVolumeOutput{}
	output.Image.ExtractImage(uploadOutput)
	return
}
