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

type ListProjectAclsInput struct {
	pagination
	Where     []string `json:"where" form:"where"`
	Namespace *string  `json:"namespace" form:"namespace"`
	_         struct{}
}

type ListProjectAclsOutput struct {
	ProjectAcls []ProjectAcl `json:"projectAcls"`
	Total       int          `json:"total"`
	_           struct{}
}

func ListProjectAcls(c *gin.Context) {
	var (
		input      = &ListProjectAclsInput{}
		output     = &ListProjectAclsOutput{}
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

	listProjectAclsInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     input.Where,
		Namespace: input.Namespace,
	}
	listProjectAclsOutput, err := vrm.ListProjectAcls(listProjectAclsInput, c)
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
			zap.String(cnt.Controller, "vrm.ListProjectAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listProjectAclsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ProjectAcls = []ProjectAcl{}
	output.Total = int(listProjectAclsOutput.Count)
	for _, data := range listProjectAclsOutput.Data {
		m := ProjectAcl{}
		m.ExtractByProto(c, data)
		output.ProjectAcls = append(output.ProjectAcls, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
