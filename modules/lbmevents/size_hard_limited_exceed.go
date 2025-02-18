package lbmevents

import "pegasus-cloud.com/aes/toolkits/littlebell"

type (
	SizeHardLimitedExceedEvent struct{ littlebell.MessageStruct }

	SizeHardLimitedExceed struct {
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

func (e *SizeHardLimitedExceedEvent) Name() string {
	return "VRM_SIZE_HARD_LIMITED_EXCEED"
}
