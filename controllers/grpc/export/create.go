package export

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

func (m *Method) CreateExport(ctx context.Context, input *pb.ExportInfo) (output *pb.ExportInfo, err error) {
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

	createInput := &storCom.CreateExportInput{
		Export: tables.Export{
			ID:             input.ID,
			RepositoryID:   input.RepositoryID,
			RepositoryName: input.RepositoryName,
			TagID:          input.TagID,
			TagName:        input.TagName,
			Type:           input.Type,
			SnapshotID:     input.SnapshotID,
			SnapshotStatus: input.SnapshotStatus,
			VolumeID:       input.VolumeID,
			VolumeStatus:   input.VolumeStatus,
			ImageID:        input.ImageID,
			ImageStatus:    input.ImageStatus,
			Filepath:       input.Filepath,
			Status:         input.Status,
			Creator:        input.Creator,
			ProjectID:      input.ProjectID,
			Namespace:      input.Namespace,
		},
	}

	createOutput, err := storages.Use().CreateExport(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageExportNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCExportNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "stor.Use().CreateExport()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &createOutput.Export)
	return
}
