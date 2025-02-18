package tasks

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/lbmevents"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/littlebell"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

/*
CountHardLimitedExceed ...

errors:
- 18000000(internal server error)
*/
func CountHardLimitedExceed() (err error) {
	var (
		funcName  = tkUtils.NameOfFunction().Name()
		requestID = tracer.EmptyRequestID
	)
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(requestID), funcName)
	defer f(tracer.Attributes{
		"err": &err,
	})

	listProjectsInput := &pb.ListInput{}
	listProjectsOutput, err := vrm.ListProjects(listProjectsInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Task, "vrm.ListProjects(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listProjectsInput),
		).Warn(err.Error())
		err = tkErr.New(cnt.TaskInternalServerErr, err)
		return
	}

	for _, project := range listProjectsOutput.Data {
		projectQuota, _err := getProjectQuota(ctx, project.ID)
		if _err != nil {
			zap.L().With(
				zap.String(cnt.Task, "getProjectQuota(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("project-id", project.ID),
			).Warn(_err.Error())
			continue
		}

		if projectQuota.LimitCount == UNLIMITED || projectQuota.CurrentCount < projectQuota.LimitCount {
			continue
		}

		// 額度已滿 觸發警告 發送給 LBM
		event := &lbmevents.CountHardLimitedExceedEvent{}
		event.With(lbmevents.CountHardLimitedExceed{
			AvailabilityDistrict: mviper.GetString("VirtualRegistryManagement.scopes.availability_district"),
			ProjectID:            projectQuota.ID,
			ProjectName:          projectQuota.Name,
			Usage:                int64(projectQuota.CurrentCount),
			Limit:                int64(projectQuota.LimitCount),
		})
		littlebell.Publish(ctx, &littlebell.LittleBellPublishInput{
			Target: projectQuota.ID,
			Event:  event,
		})
	}
	return
}
