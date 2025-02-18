package openstack

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) CreateVolume(ctx context.Context, input *pb.VolumeInfo) (output *pb.VolumeInfo, err error) {
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
	// create volume
	createVolumeInput := &opstkCom.CreateVolumeInput{
		SnapshotID:  input.SnapshotID,
		Size:        int(input.Size / 1024 / 1024 / 1024),
		Name:        input.Name,
		Description: fmt.Sprintf("create from snapshot %s", input.SnapshotID),
	}
	createVolumeOutput, createVolumeErr := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).CreateVolume(ctx, createVolumeInput)
	if createVolumeErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "openstack.Namespace().Cinder().CreateVolume(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createVolumeInput),
		).Error(createVolumeErr.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(createVolumeErr)
		return
	}

	output = &pb.VolumeInfo{
		ID:          createVolumeOutput.Volume.ID,
		SnapshotID:  createVolumeOutput.Volume.SnapshotID,
		Name:        createVolumeOutput.Volume.Name,
		Status:      createVolumeOutput.Volume.Status.String(),
		Description: createVolumeOutput.Volume.Description,
		Size:        uint64(createVolumeOutput.Volume.Size),
		UserID:      input.UserID,
		ProjectID:   input.ProjectID,
		Namespace:   input.Namespace,
		CreatedAt:   createVolumeOutput.Volume.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   createVolumeOutput.Volume.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return
}
