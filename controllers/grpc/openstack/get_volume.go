package openstack

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"
	"time"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

func (m *Method) GetVolume(ctx context.Context, input *pb.GetOpenstackInput) (output *pb.VolumeInfo, err error) {
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

	getVolumeInput := &opstkCom.GetVolumeInput{
		ID: input.ID,
	}
	getOutput, getVolumeErr := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).GetVolume(ctx, getVolumeInput)
	if getVolumeErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "openstack.Namespace().Cinder().GetVolume(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getVolumeInput),
		).Error(getVolumeErr.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getVolumeErr)
		return
	}

	output = &pb.VolumeInfo{
		ID:          getOutput.Volume.ID,
		SnapshotID:  getOutput.Volume.SnapshotID,
		Name:        getOutput.Volume.Name,
		Status:      getOutput.Volume.Status.String(),
		Description: getOutput.Volume.Description,
		Size:        uint64(getOutput.Volume.Size),
		ProjectID:   input.ProjectID,
		Namespace:   input.Namespace,
		CreatedAt:   getOutput.Volume.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   getOutput.Volume.UpdatedAt.UTC().Format(time.RFC3339),
	}
	return
}
