package user

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

type ListExportsInput struct {
	pagination
	Where []string `json:"where" form:"where"`

	ProjectID string `json:"-"`
	_         struct{}
}

type ListExportsOutput struct {
	Exports []Export `json:"exports"`
	Total   int      `json:"total"`
	_       struct{}
}

func ListExports(c *gin.Context) {
	var (
		input = &ListExportsInput{
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output       = &ListExportsOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusOK
		supportRoles     = map[string]bool{
			cnt.TenantOwner.String(): true,
			cnt.TenantAdmin.String(): true,
		}
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

	// tenant-owner and tenant-admin allow to get all in project
	if role, userID := c.GetString(cnt.CtxTenantRole), c.GetString(cnt.CtxUserID); !supportRoles[role] && userID != "" {
		input.Where = append(input.Where, "creator="+userID)
	}

	namespace := c.GetString(cnt.CtxNamespace)
	listExportsInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     append(input.Where, "project-id="+input.ProjectID),
		Namespace: &namespace,
	}
	listExportsOutput, err := vrm.ListExports(listExportsInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				statusCode = http.StatusBadRequest
				if v, exist := e.Get("field"); exist {
					err = tkErr.New(cnt.UserAPIQueryNotSupportErr, "where", v)
				} else {
					err = tkErr.New(cnt.UserAPIIllegalWhereQueryFormatErr)
				}
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.ListExports(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listExportsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Exports = []Export{}
	output.Total = int(listExportsOutput.Count)
	for _, data := range listExportsOutput.Data {
		m := Export{}
		m.ExtractByProto(c, data)
		output.Exports = append(output.Exports, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
