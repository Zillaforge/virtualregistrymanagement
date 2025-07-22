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
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

func UnmarshalDeleteMembership(input *ecCom.Data) (output *iamPB.MemUserProjectInput) {
	output = &iamPB.MemUserProjectInput{}
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

func DeleteMembership(ctx context.Context, input *iamPB.MemUserProjectInput) {
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

	listInput := &pb.ListRegistriesInput{
		Where: []string{
			"allow-user-id=" + input.UserID,
			"project-id=" + input.ProjectID,
		},
	}
	listOutput, err := vrm.ListRegistries(listInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.EventConsume, "vrm.ListRegistries(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(err.Error())
		return
	}

	for _, tag := range listOutput.Data {
		deleteInput := &pb.DeleteInput{
			Where: []string{"TagID=" + tag.TagID},
		}
		if _, err = vrm.DeleteMemberAcl(deleteInput, ctx); err != nil {
			zap.L().With(
				zap.String(cnt.EventConsume, "vrm.DeleteMemberAcl(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", deleteInput),
			).Error(err.Error())
			return
		}
	}
}
