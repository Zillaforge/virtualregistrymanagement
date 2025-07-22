package tasks

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"

	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/scheduler"
)

const (
	ignoreTasks = "off"
)

type Task struct {
	/*
		Explanation:
			CronKey => The key in yaml, basically is the name of scheduler task.
			JobFunc => The main task corresponding to "CronKey".
			JobName => For log use.
			Execute => Does the task need to execute ahead of the scheduler trigger it.
	*/
	CronKey string
	JobFunc func() error
	JobName string
	Execute bool
}

var TaskList = []Task{
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.sync_projects.cron_expression",
		JobFunc: SyncProjects,
		JobName: "SyncProjects",
		Execute: true,
	},
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.sync_tags.cron_expression",
		JobFunc: SyncTags,
		JobName: "SyncTags",
		Execute: true,
	},
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.sync_exports.cron_expression",
		JobFunc: SyncExports,
		JobName: "SyncExports",
		Execute: true,
	},
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.image_count.cron_expression",
		JobFunc: CalculateProjectImages,
		JobName: "CalculateProjectImages",
		Execute: true,
	},
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.image_size.cron_expression",
		JobFunc: CalculateProjectSize,
		JobName: "CalculateProjectSize",
		Execute: true,
	},
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.size_hard_limited_exceed.cron_expression",
		JobFunc: SizeHardLimitedExceed,
		JobName: "SizeHardLimitedExceed",
		Execute: true,
	},
	{
		CronKey: "VirtualRegistryManagementScheduler.tasks.count_hard_limited_exceed.cron_expression",
		JobFunc: CountHardLimitedExceed,
		JobName: "CountHardLimitedExceed",
		Execute: true,
	},
}

func InitSchedulerTasks() {
	for _, taskInfo := range TaskList {

		// If cron expression gets "off" meaning don't need to set up that schedule job
		if cronExpression := mviper.GetString(taskInfo.CronKey); cronExpression == ignoreTasks {
			zap.L().With(zap.String(cnt.Task, taskInfo.JobName)).Info("Task Ignore")
			continue
		}

		task := scheduler.CreateSchedulerV2(mviper.GetString(taskInfo.CronKey))
		if err := task.Time().Do(taskInfo.JobFunc); err != nil {
			panic(fmt.Sprintf(
				"start %s(cronExp: %s) failed: %s",
				taskInfo.JobName,
				mviper.GetString(taskInfo.CronKey),
				err.Error()))
		}
		task.Start()
		// Execute the job ahead of scheduler triggered
		if taskInfo.Execute {
			taskInfo.JobFunc()
		}
	}
}
