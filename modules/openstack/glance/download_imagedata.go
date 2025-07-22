package glance

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"
	"errors"
	"io"
	"os"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (g *Glance) DownloadImageData(ctx context.Context, input *common.DownloadImageDataInput) (output *common.DownloadImageDataOutput, err error) {
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

	data, err := imagedata.Download(g.sc, input.ID).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "imagedata.Download(...).ExtractErr()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.String("image-id", input.ID),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}
	defer data.Close()

	if _, err = os.Stat(input.Filepath); !errors.Is(err, os.ErrNotExist) {
		zap.L().With(
			zap.String(cnt.Module, "os.Stat()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.String("image-id", input.ID),
			zap.String("filepath", input.Filepath),
		).Info(cnt.OpenstackFileIsExistErrMsg)
		err = tkErr.New(cnt.OpenstackFileIsExistErr).WithInner(err)
		return
	}

	file, err := os.Create(input.Filepath)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "os.Create()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.String("image-id", input.ID),
			zap.String("filepath", input.Filepath),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackCreateFileFailedErr).WithInner(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "io.Copy()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", g.namespace),
			zap.String("image-id", input.ID),
			zap.String("filepath", input.Filepath),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackImageToFileFailedErr).WithInner(err)
		return
	}

	output = &common.DownloadImageDataOutput{}
	return
}
