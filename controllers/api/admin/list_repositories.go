package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type ListRepositoriesInput struct {
	pagination
	Where     []string `json:"where" form:"where"`
	Namespace *string  `json:"namespace" form:"namespace"`
	_         struct{}
}

type ListRepositoriesOutput struct {
	Repositories []Repository `json:"repositories"`
	Total        int          `json:"total"`
	_            struct{}
}

func ListRepositories(c *gin.Context) {
	var (
		input      = &ListRepositoriesInput{}
		output     = &ListRepositoriesOutput{}
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

	listRepositoriesInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     input.Where,
		Namespace: input.Namespace,
	}
	listRepositoriesOutput, err := vrm.ListRepositories(listRepositoriesInput, c)
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
			zap.String(cnt.Controller, "vrm.ListRepositories(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRepositoriesInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Repositories = []Repository{}
	output.Total = int(listRepositoriesOutput.Count)
	for _, data := range listRepositoriesOutput.Data {
		m := Repository{}
		m.ExtractByProto(c, data)
		output.Repositories = append(output.Repositories, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
