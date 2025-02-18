package cinder

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (c *Cinder) GetVolume(ctx context.Context, input *common.GetVolumeInput) (output *common.GetVolumeOutput, err error) {
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
			zap.String(cnt.Module, "c.checkConnection()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
		).Error(err.Error())
		return
	}

	getOutput, err := volumes.Get(c.sc, input.ID).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "volumes.Get(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.Any("volume-id", input.ID),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.GetVolumeOutput{}
	output.Volume.ExtractVolume(getOutput)
	return
}
