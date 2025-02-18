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
	iamPB "pegasus-cloud.com/aes/pegasusiamclient/pb"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

func UnmarshalDeleteUser(input *ecCom.Data) (output *iamPB.UserID) {
	output = &iamPB.UserID{}
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

func DeleteUser(ctx context.Context, input *iamPB.UserID) {
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

	deleteInput := &pb.DeleteInput{
		Where: []string{"UserID=" + input.ID},
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
