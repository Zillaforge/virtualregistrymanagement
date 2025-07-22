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

func (m *Method) CreateSnapshot(ctx context.Context, input *pb.SnapshotInfo) (output *pb.SnapshotInfo, err error) {
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

	createSnapshotInput := &opstkCom.CreateSnapshotInput{
		VolumeID: input.VolumeID,
		Name:     input.Name,
		Force:    true,

		Creator: &input.UserID,
	}
	createSnapshotOutput, _err := openstack.Namespace(input.Namespace).Cinder(input.ProjectID).CreateSnapshot(ctx, createSnapshotInput)
	if _err != nil {
		if e, ok := tkErr.IsError(_err); ok {
			switch e.Code() {
			case cnt.OpenstackExceedAllowedQuotaErr.Code():
				err = tkErr.New(cCnt.GRPCExceedAllowedQuotaErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "openstack.Namespace().Cinder().CreateSnapshot()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", input.Namespace),
			zap.String("project-id", input.ProjectID),
			zap.Any("input", createSnapshotInput),
		).Error(_err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
		return
	}

	output = &pb.SnapshotInfo{
		ID:          createSnapshotOutput.Snapshot.ID,
		VolumeID:    createSnapshotOutput.Snapshot.VolumeID,
		Name:        createSnapshotOutput.Snapshot.Name,
		Status:      createSnapshotOutput.Snapshot.Status.String(),
		Description: createSnapshotOutput.Snapshot.Description,
		Size:        uint64(createSnapshotOutput.Snapshot.FullSizeBytes),
		UserID:      input.UserID,
		ProjectID:   input.ProjectID,
		Namespace:   input.Namespace,
		CreatedAt:   createSnapshotOutput.Snapshot.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   createSnapshotOutput.Snapshot.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return
}
