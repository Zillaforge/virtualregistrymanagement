package cinder

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (c *Cinder) UpdateSnapshotMetadata(ctx context.Context, input *common.UpdateSnapshotMetadataInput) (output *common.UpdateSnapshotMetadataOutput, err error) {
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

	updateMetadataInput := input.UpdateOpts(c.namespace, c.projectID)
	updateMetadataOutput, err := snapshots.UpdateMetadata(c.sc, input.ID, updateMetadataInput).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "snapshots.Delete(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.String("image-id", input.ID),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.UpdateSnapshotMetadataOutput{}
	output.Snapshot.ExtractSnapshot(updateMetadataOutput)
	return
}
