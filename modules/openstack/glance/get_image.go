package glance

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (g *Glance) GetImage(ctx context.Context, input *common.GetImageInput) (output *common.GetImageOutput, err error) {
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

	if err = g.checkConnection(); err != nil {
		zap.L().With(
			zap.String(cnt.Module, "g.checkConnection()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
		).Error(err.Error())
		return
	}

	getOutput, err := images.Get(g.sc, input.ID).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "images.Get(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.Any("image-id", input.ID),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.GetImageOutput{}
	output.Image.ExtractImage(getOutput)
	return
}
