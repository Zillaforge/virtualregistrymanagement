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

func (c *Cinder) ListSnapshots(ctx context.Context, input *common.ListSnapshotsInput) (output *common.ListSnapshotsOutput, err error) {
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

	listInput := input.ListOpts(c.namespace)
	page, err := snapshots.List(c.sc, listInput).AllPages()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "snapshots.List(...).AllPages()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.Any("input", listInput),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}
	listOutput, err := snapshots.ExtractSnapshots(page)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "snapshots.ExtractSnapshots(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", c.namespace),
			zap.Any("input", listInput),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.ListSnapshotsOutput{}
	for _, snapshot := range listOutput {
		s := &common.SnapshotInfo{}
		output.Snapshots = append(output.Snapshots, s.ExtractSnapshot(&snapshot))
	}
	return
}
