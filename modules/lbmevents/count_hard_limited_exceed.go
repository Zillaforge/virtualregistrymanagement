package lbmevents

import "pegasus-cloud.com/aes/toolkits/littlebell"

type (
	CountHardLimitedExceedEvent struct{ littlebell.MessageStruct }

	CountHardLimitedExceed struct {
		// @message AD
		// AD
		AvailabilityDistrict string `json:"ad"`
		// @message projectID
		// Project ID
		ProjectID string `json:"projectID"`
		// @message projectName
		// Project Name
		ProjectName string `json:"projectName"`
		// @message limit
		// Limit
		Limit int64 `json:"limit"`
		// @message usage
		// Usage
		Usage int64 `json:"usage"`

		_ struct{}
	}
)

func (e *CountHardLimitedExceedEvent) Name() string {
	return "VRM_COUNT_HARD_LIMITED_EXCEED"
}
