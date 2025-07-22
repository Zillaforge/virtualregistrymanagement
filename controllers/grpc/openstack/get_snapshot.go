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

func (m *Method) GetSnapshot(ctx context.Context, input *pb.GetOpenstackInput) (output *pb.SnapshotInfo, err error) {
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

	getSnapshotInput := &opstkCom.GetSnapshotInput{
		ID: input.ID,
	}
	getOutput, getErr := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).GetSnapshot(ctx, getSnapshotInput)
	if getErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "openstack.Namespace().Cinder().GetSnapshot(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getSnapshotInput),
		).Error(getErr.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getErr)
		return
	}

	output = &pb.SnapshotInfo{
		ID:          getOutput.Snapshot.ID,
		VolumeID:    getOutput.Snapshot.VolumeID,
		Name:        getOutput.Snapshot.Name,
		Status:      getOutput.Snapshot.Status.String(),
		Description: getOutput.Snapshot.Description,
		Size:        uint64(getOutput.Snapshot.FullSizeBytes),
		ProjectID:   input.ProjectID,
		Namespace:   input.Namespace,
		CreatedAt:   getOutput.Snapshot.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   getOutput.Snapshot.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return
}
