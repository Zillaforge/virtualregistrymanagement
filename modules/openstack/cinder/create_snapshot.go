package cinder

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (c Cinder) CreateSnapshot(ctx context.Context, input *common.CreateSnapshotInput) (output *common.CreateSnapshotOutput, err error) {
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

	createInput := snapshots.CreateOpts{
		VolumeID:    input.VolumeID,
		Name:        common.AddPrefixName(input.Name),
		Description: input.Description,
		Force:       input.Force,
		Metadata:    input.Tag(c.namespace, c.projectID),
	}
	createOutput, err := snapshots.Create(c.sc, createInput).Extract()
	if err != nil {
		log := zap.L().With(
			zap.String(cnt.Module, "snapshots.Create(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.Any("input", createInput),
		)

		if e, ok := common.IsKnownError(err); ok {
			log.Warn(err.Error())
			err = e
			return
		}
		log.Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.CreateSnapshotOutput{}
	output.Snapshot.ExtractSnapshot(createOutput)
	return
}
