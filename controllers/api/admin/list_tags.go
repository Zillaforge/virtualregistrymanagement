package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type ListTagsInput struct {
	pagination
	Where     []string `json:"where" form:"where"`
	Namespace *string  `json:"namespace" form:"namespace"`
	_         struct{}
}

type ListTagsOutput struct {
	Tags  []Tag `json:"tags"`
	Total int   `json:"total"`
	_     struct{}
}

func ListTags(c *gin.Context) {
	var (
		input      = &ListTagsInput{}
		output     = &ListTagsOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindQuery(input); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindQuery(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	listTagsInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     input.Where,
		Namespace: input.Namespace,
	}
	listTagsOutput, err := vrm.ListTags(listTagsInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				statusCode = http.StatusBadRequest
				if v, exist := e.Get("field"); exist {
					err = tkErr.New(cnt.AdminAPIQueryNotSupportErr, "where", v)
				} else {
					err = tkErr.New(cnt.AdminAPIIllegalWhereQueryFormatErr)
				}
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.ListTags(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listTagsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Tags = []Tag{}
	output.Total = int(listTagsOutput.Count)
	for _, data := range listTagsOutput.Data {
		m := Tag{}
		m.ExtractByProto(c, data)
		output.Tags = append(output.Tags, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
