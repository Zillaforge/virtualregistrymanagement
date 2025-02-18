package openstack

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) DeleteSnapshot(ctx context.Context, input *pb.SnapshotInfo) (output *pb.DeleteOutput, err error) {
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

	deleteSnapshotInput := &opstkCom.DeleteSnapshotInput{
		ID: input.ID,
	}
	_, err = openstack.Namespace(input.Namespace).Cinder(input.ProjectID).DeleteSnapshot(ctx, deleteSnapshotInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.OpenstackExceedAllowedQuotaErr.Code():
				err = tkErr.New(cCnt.GRPCExceedAllowedQuotaErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "openstack.Namespace().Cinder().DeleteSnapshot()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", input.Namespace),
			zap.String("project-id", input.ProjectID),
			zap.Any("input", deleteSnapshotInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = &pb.DeleteOutput{
		ID: []string{input.ID},
	}
	return
}
