package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/modules/openstack"
	opstkCom "VirtualRegistryManagement/modules/openstack/common"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"
	"encoding/json"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

// 若 reference_target 有值，則從 openstack 取得其相關資訊並創建 tag
// 反之，創建 tag 後再到 openstack 建立 image/snapshot 相關資訊
func (m *Method) CreateTag(ctx context.Context, input *pb.CreateTagInput) (output *pb.TagDetail, err error) {
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

	getInput := &pb.GetInput{
		ID: input.Tag.RepositoryID,
	}
	getOutput, err := vrm.GetRepository(getInput, ctx)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCRepositoryNotFoundErr.Code():
				err = e
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.GRPC, "vrm.GetRepository()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	var (
		namespace      = getOutput.Repository.Namespace
		projectID      = getOutput.Repository.ProjectID
		repositoryName = getOutput.Repository.Name
	)

	// create tag and openstack resource
	createInput := &storCom.CreateTagInput{
		Tag: tables.Tag{
			ID:           input.Tag.ID,
			Name:         input.Tag.Name,
			RepositoryID: input.Tag.RepositoryID,
			Type:         input.Tag.Type,
			Size:         input.Tag.Size,
			Status:       input.Tag.Status,
			Extra:        input.Tag.Extra,
		},
	}

	var image *opstkCom.ImageInfo
	// get information from openstack, and create tag
	if input.Tag.ReferenceTarget != "" && input.Tag.Type == grpc.TagTypeImage {
		getImageInput := &opstkCom.GetImageInput{
			ID: input.Tag.ReferenceTarget,
		}
		getImageOutput, _err := openstack.Namespace(namespace).Glance(projectID).GetImage(ctx, getImageInput)
		if _err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Glance().GetImage()"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", namespace),
				zap.String("project-id", projectID),
				zap.Any("input", getImageInput),
			).Error(_err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
			return
		}

		// validate properties
		for _, key := range []string{"defaultUser", "distribution"} {
			if _, ok := getImageOutput.Image.Properties[key]; !ok {
				zap.L().With(
					zap.String(cnt.GRPC, "openstack.Namespace().Glance().GetImage()"),
					zap.String(cnt.RequestID, requestID),
					zap.String("namespace", namespace),
					zap.String("image-id", getImageOutput.Image.ID),
					zap.Any("input", getImageOutput.Image.Properties),
				).Info("image loss properties")
				err = tkErr.New(cCnt.GRPCImportImageLossPropertiesErr)
				return
			}
		}

		image = &getImageOutput.Image

		name := image.Name
		if input.Tag.Name != "" {
			name = input.Tag.Name
		}
		extra, _ := json.Marshal(&image.Properties)
		createInput = &storCom.CreateTagInput{
			Tag: tables.Tag{
				ID:              input.Tag.ID,
				Name:            name,
				RepositoryID:    input.Tag.RepositoryID,
				ReferenceTarget: input.Tag.ReferenceTarget,
				Type:            grpc.TagTypeImage,
				Size:            uint64(image.SizeBytes),
				Status:          image.Status.String(),
				Extra:           extra,
			},
		}
	}

	extra := map[string]interface{}{}
	json.Unmarshal(createInput.Tag.Extra, &extra)
	extra["diskFormat"] = "raw"
	if input.Tag.Type == grpc.TagTypeImage && input.Image != nil {
		extra["diskFormat"] = input.Image.DiskFormat
	}
	byteExtra, _ := json.Marshal(extra)
	createInput.Tag.Extra = byteExtra

	createOutput, err := storages.Use().CreateTag(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageTagExistErr.Code():
				err = tkErr.New(cCnt.GRPCTagExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateTag()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// update information to openstack image tags
	if input.Tag.ReferenceTarget != "" && input.Tag.Type == grpc.TagTypeImage {
		getTagInput := &storCom.GetTagInput{ID: createOutput.Tag.ID}
		getTagOutput, getTagErr := storages.Use().GetTag(ctx, getTagInput)
		if getTagErr != nil {
			if e, ok := tkErr.IsError(getTagErr); ok {
				switch e.Code() {
				case cnt.StorageTagExistErr.Code():
					err = tkErr.New(cCnt.GRPCTagExistErr)
					return
				}
			}
			zap.L().With(
				zap.String(cnt.GRPC, "storages.Use().GetTag()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", createInput),
			).Error(getTagErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(getTagErr)
			return
		}

		// add vrm tags to the existing ones
		updateImageInput := &opstkCom.UpdateImageInput{
			ID:   image.ID,
			Tags: &image.Tags,

			Creator:      &getTagOutput.Tag.Repository.Creator,
			RepositoryID: &getTagOutput.Tag.Repository.ID,
			TagID:        &getTagOutput.Tag.ID,
		}
		_, updateImageErr := openstack.Namespace(namespace).Glance(projectID).UpdateImage(ctx, updateImageInput)
		if updateImageErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace().Glance().UpdateImage()"),
				zap.String(cnt.RequestID, requestID),
				zap.String("namespace", namespace),
				zap.String("project-id", projectID),
				zap.Any("input", updateImageInput),
			).Error(updateImageErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(updateImageErr)
			return
		}
	}

	// resource is not created
	tagName := repositoryName + "-" + input.Tag.Name + "-" + namespace
	if createOutput.Tag.ReferenceTarget == "" {
		var (
			resourceID string
			status     string
		)
		switch input.Tag.Type {
		case grpc.TagTypeImage:
			createImageInput := &opstkCom.CreateImageInput{
				Name:            tagName,
				DiskFormat:      input.Image.DiskFormat,
				ContainerFormat: input.Image.ContainerFormat,
				Visibility:      opstkCom.ImageVisibilityPrivate,

				TagID:        &createOutput.Tag.ID,
				RepositoryID: &getOutput.Repository.ID,
				Creator:      &getOutput.Repository.Creator,
			}
			createImageOutput, _err := openstack.Namespace(namespace).Glance(projectID).CreateImage(ctx, createImageInput)
			if _err != nil {
				m.DeleteTag(ctx, &pb.DeleteInput{Where: []string{"ID=" + createOutput.Tag.ID}})

				if e, ok := tkErr.IsError(_err); ok {
					switch e.Code() {
					case cnt.OpenstackExceedAllowedQuotaErr.Code():
						err = tkErr.New(cCnt.GRPCExceedAllowedQuotaErr)
						return
					}
				}
				zap.L().With(
					zap.String(cnt.GRPC, "openstack.Namespace().Glance().CreateImage()"),
					zap.String(cnt.RequestID, requestID),
					zap.String("namespace", namespace),
					zap.String("project-id", projectID),
					zap.Any("input", createImageInput),
				).Error(_err.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
				return
			}
			resourceID = createImageOutput.Image.ID
			status = createImageOutput.Image.Status.String()

		case grpc.TagTypeVolumeSnapshot:
			createSnapshotInput := &opstkCom.CreateSnapshotInput{
				VolumeID: input.Snapshot.VolumeID,
				Name:     tagName,
				Force:    true,

				TagID:        &createOutput.Tag.ID,
				RepositoryID: &getOutput.Repository.ID,
				Creator:      &getOutput.Repository.Creator,
			}
			createSnapshotOutput, _err := openstack.Namespace(namespace).Cinder(projectID).CreateSnapshot(ctx, createSnapshotInput)
			if _err != nil {
				m.DeleteTag(ctx, &pb.DeleteInput{Where: []string{"ID=" + createOutput.Tag.ID}})

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
					zap.String("namespace", namespace),
					zap.String("project-id", projectID),
					zap.Any("input", createSnapshotInput),
				).Error(_err.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
				return
			}
			resourceID = createSnapshotOutput.Snapshot.ID
			status = createSnapshotOutput.Snapshot.Status.String()
		}

		updateInput := &pb.UpdateTagInput{
			ID:              createOutput.Tag.ID,
			ReferenceTarget: &resourceID,
			Status:          &status,
		}
		_, _err := vrm.UpdateTag(updateInput, ctx)
		if _err != nil {
			if _, ok := tkErr.IsError(_err); ok {
				err = tkErr.New(cCnt.GRPCInternalServerErr)
				return
			}
			zap.L().With(
				zap.String(cnt.GRPC, "vrm.UpdateTag()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", updateInput),
			).Error(_err.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
			return
		}
	}

	output = m.storage2pb(ctx, &createOutput.Tag)
	return
}
