package glance

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"
	"os"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (g *Glance) UploadImageData(ctx context.Context, input *common.UploadImageDataInput) (output *common.UploadImageDataOutput, err error) {
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
		err = tkErr.New(cnt.OpenstackUploadFileNotFoundErr).WithInner(err)
		return
	}

	data, err := os.Open(input.Filepath)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "os.Open()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.String("filepath", input.Filepath),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackUploadFileNotFoundErr).WithInner(err)
		return
	}
	defer data.Close()

	err = imagedata.Upload(g.sc, input.ID, data).ExtractErr()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "imagedata.Upload(...).ExtractErr()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.String("image-id", input.ID),
			zap.String("filepath", input.Filepath),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}

	output = &common.UploadImageDataOutput{}
	return
}
