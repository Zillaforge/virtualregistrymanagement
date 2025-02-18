package tasks

import (
	cnt "VirtualRegistryManagement/constants"
	"time"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/meteringtoolkits/metering"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"

	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

const (
	_calculateProjectImages = "CalculateProjectImages"
)

func CalculateProjectImages() (err error) {
	current := time.Now()
	var (
		funcName  = tkUtils.NameOfFunction().Name()
		requestID = tracer.EmptyRequestID
	)
	zap.L().With(zap.String(cnt.Task, _calculateProjectImages)).Debug("Task Work")

	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(requestID), funcName)
	defer f(tracer.Attributes{
		"err": &err,
	})

	// ListAllProjects
	listProjectsInput := &pb.ListInput{
		Limit:  -1,
		Offset: 0,
	}
	listProjectsOutput, err := vrm.ListProjects(listProjectsInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Task, "vrm.ListProjects(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listProjectsInput),
		).Error(err.Error())
		return
	}

	// For every project need to list all repository to compute
	for _, project := range listProjectsOutput.Data {
		where := []string{"project-id=" + project.ID}
		listRepositoriesInput := &pb.ListNamespaceInput{
			Limit:  -1,
			Offset: 0,
			Where:  where,
		}
		listRepositoriesOutput, _err := vrm.ListRepositories(listRepositoriesInput, ctx)
		if _err != nil {
			zap.L().With(
				zap.String(cnt.Task, "vrm.ListRepositories(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", listRepositoriesInput),
			).Error(_err.Error())
			_err = tkErr.New(cnt.TaskInternalServerErr, _err)
			continue
		}

		// Compute total size by walk through all repositories
		var currentCount uint64 = 0

		for _, repository := range listRepositoriesOutput.Data {
			currentCount += uint64(len(repository.Tags))
		}

		// If the project has 0 usage counts, do not publish the message or it will throw an error back. (body.value is required)
		if currentCount != 0 {
			publishInput := &metering.PublishInput{
				Exchange:   mviper.GetString("VirtualRegistryManagementScheduler.tasks.image_count.metering_service.exchange"),
				RoutingKey: mviper.GetString("VirtualRegistryManagementScheduler.tasks.image_count.metering_service.routing_key"),
				Body: metering.PublishBody{
					AvailabilityDistrict: mviper.GetString("VirtualRegistryManagement.scopes.availability_district"),
					Service:              cnt.Name,
					ProjectID:            project.ID,
					UserID:               "99999999-9999-9999-999999999999",
					ResourceType:         "vrm-image-count",
					ResourceID:           project.ID,
					ResourceName:         "vrm-image-count",
					RecordTime:           current,
					StartTime:            current,
					EndTime:              current,
					Value:                currentCount,
					Unit:                 metering.PicUnit,
				},
			}

			if publishErr := metering.Publish(ctx, publishInput); publishErr != nil {
				zap.L().With(
					zap.String(cnt.Task, "metering.Publish(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", publishInput),
				).Error(publishErr.Error())
				continue
			}
		}
	}

	return
}
