package glance

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (g *Glance) CreateImage(ctx context.Context, input *common.CreateImageInput) (output *common.CreateImageOutput, err error) {
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

	createInput := images.CreateOpts{
		Name:            common.AddPrefixName(input.Name),
		DiskFormat:      input.DiskFormat,
		ContainerFormat: input.ContainerFormat,
		Visibility:      input.Visibility.Convert(),
		Tags:            input.Tag(g.namespace, g.projectID),
	}

	createOutput, err := images.Create(g.sc, createInput).Extract()
	if err != nil {
		log := zap.L().With(
			zap.String(cnt.Module, "images.Create(...).Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
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

	output = &common.CreateImageOutput{}
	output.Image.ExtractImage(createOutput)
	return
}
