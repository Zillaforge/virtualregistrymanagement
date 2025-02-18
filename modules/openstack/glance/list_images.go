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

func (g *Glance) ListImages(ctx context.Context, input *common.ListImagesInput) (output *common.ListImagesOutput, err error) {
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

	listInput := input.ListOpts(g.namespace)
	page, err := images.List(g.sc, listInput).AllPages()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "images.List(...).AllPages()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.Any("input", listInput),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}
	listOutput, err := images.ExtractImages(page)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "images.ExtractImages(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.Any("input", listInput),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.ListImagesOutput{}
	for _, image := range listOutput {
		img := &common.ImageInfo{}
		output.Images = append(output.Images, img.ExtractImage(&image))
	}
	return
}
