package iamconsumer

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/logger"
	ecCom "VirtualRegistryManagement/modules/eventconsume/common"
	"VirtualRegistryManagement/utility"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	iamPB "github.com/Zillaforge/pegasusiamclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

func UnmarshalCreateProject(input *ecCom.Data) (output *iamPB.ProjectID) {
	output = &iamPB.ProjectID{}
	if input.Request != nil {
		switch v := input.Request.(type) {
		case string:
			decodedData, _ := base64.StdEncoding.DecodeString(v)
			json.Unmarshal(decodedData, output)
			logger.Use().Info(fmt.Sprintf("%s | %s | %s | %s",
				input.Metadata[tracer.RequestID],
				"FromIAM",
				input.Action,
				decodedData,
			))
		default:
			zap.L().Warn(fmt.Sprintf("Received the message format of %s action is invalid", input.Action))
		}
	}
	return output
}

func CreateProject(ctx context.Context, input *iamPB.ProjectID) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
		err       error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input": &input,
			"err":   &err,
		},
	)

	createInput := &pb.ProjectInfo{
		ID:             input.ID,
		LimitCount:     mviper.GetInt64("VirtualRegistryManagement.scopes.project_default_count"),
		LimitSizeBytes: mviper.GetInt64("VirtualRegistryManagement.scopes.project_default_size"),
	}
	if _, err = vrm.CreateProject(createInput, ctx); err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCProjectExistErr.Code():
				err = tkErr.New(cnt.TaskProjectExistErr, input.ID)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.EventConsume, "vrm.CreateProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		return
	}
}
