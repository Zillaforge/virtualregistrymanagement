package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/modules/openstack"
	"VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

const (
	_do_not_do_anything = "don't do anything"
	_create             = "create"
	_delete             = "delete"
)

func (m *Method) SyncTags(ctx context.Context, input *emptypb.Empty) (output *emptypb.Empty, err error) {
	var (
		funcName   = tkUtils.NameOfFunction().String()
		requestID  = utility.MustGetContextRequestID(ctx)
		projectMap = make(map[string]string)
	)
	output = empty

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &projectMap,
			"output": &output,
			"error":  &err,
		},
	)

	type info struct {
		Status string
		Size   uint64
	}

	imageMap := map[string]info{}
	snapshotMap := map[string]info{}

	labelSlice := []string{common.CreatedTag}
	labelMap := map[string]string{
		common.CreatedTag: "",
	}
	for _, namespace := range openstack.ListNamespaces() {

		listImageInput := &common.ListImagesInput{
			Tags: &labelSlice,
		}
		listImageOutput, listImageErr := openstack.Namespace(namespace).Admin().Glance().ListImages(ctx, listImageInput)
		if listImageErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Admin().Glance().ListImages(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", listImageInput),
			).Error(listImageErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(listImageErr)
		} else {
			for _, image := range listImageOutput.Images {
				imageMap[image.ID] = info{
					Status: image.Status.String(),
					Size:   uint64(image.SizeBytes),
				}
			}
		}

		listSnapshotInput := &common.ListSnapshotsInput{
			Metadata: &labelMap,
		}
		listSnapshotOutput, listSnapshotErr := openstack.Namespace(namespace).Admin().Cinder().ListSnapshots(ctx, listSnapshotInput)
		if listSnapshotErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Admin().Cinder().ListSnapshots(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", listImageInput),
			).Error(listSnapshotErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(listSnapshotErr)
		} else {
			for _, snapshot := range listSnapshotOutput.Snapshots {
				snapshotMap[snapshot.ID] = info{
					Status: snapshot.Status.String(),
					Size:   uint64(snapshot.SizeBytes),
				}
			}
		}
	}

	listTagInput := &pb.ListNamespaceInput{
		Limit:  -1,
		Offset: 0,
	}
	listTagOutput, listTagErr := m.ListTags(ctx, listTagInput)
	if listTagErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "m.ListTags(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listTagInput),
		).Error(listTagErr.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(listTagErr)
		return
	}

	for _, tag := range listTagOutput.Data {
		switch tag.Tag.Type {
		case grpc.TagTypeImage:
			if val, ok := imageMap[tag.Tag.ReferenceTarget]; ok &&
				(val.Status != tag.Tag.Status || val.Size != uint64(tag.Tag.Size)) {
				updateTagInput := &pb.UpdateTagInput{
					ID:     tag.Tag.ID,
					Size:   &val.Size,
					Status: &val.Status,
				}
				_, updateTagErr := m.UpdateTag(ctx, updateTagInput)
				if updateTagErr != nil {
					zap.L().With(
						zap.String(cnt.GRPC, "m.UpdateTag(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.Any("input", updateTagInput),
					).Error(updateTagErr.Error())
					err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(updateTagErr)
				}
			}
		case grpc.TagTypeVolumeSnapshot:
			if val, ok := snapshotMap[tag.Tag.ReferenceTarget]; ok &&
				(val.Status != tag.Tag.Status || val.Size != uint64(tag.Tag.Size)) {
				updateTagInput := &pb.UpdateTagInput{
					ID:     tag.Tag.ID,
					Size:   &val.Size,
					Status: &val.Status,
				}
				_, updateTagErr := m.UpdateTag(ctx, updateTagInput)
				if updateTagErr != nil {
					zap.L().With(
						zap.String(cnt.GRPC, "m.UpdateTag(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.Any("input", updateTagInput),
					).Error(updateTagErr.Error())
					err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(updateTagErr)
				}
			}
		}
	}
	return
}
