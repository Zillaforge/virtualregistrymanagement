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

func (c Cinder) CreateVolume(ctx context.Context, input *common.CreateVolumeInput) (output *common.CreateVolumeOutput, err error) {
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

	createInput := volumes.CreateOpts{
		SnapshotID:  input.SnapshotID,
		Size:        input.Size,
		Name:        input.Name,
		Description: input.Description,
		Metadata:    input.Metadata,
	}
	createOutput, err := volumes.Create(c.sc, createInput).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "volumes.Create(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.CreateVolumeOutput{}
	output.Volume.ExtractVolume(createOutput)
	return
}
