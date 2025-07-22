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

type ListMemberAclsInput struct {
	pagination
	Where []string `json:"where" form:"where"`
	_     struct{}
}

type ListMemberAclsOutput struct {
	MemberAcls []MemberAcl `json:"memberAcls"`
	Total      int         `json:"total"`
	_          struct{}
}

func ListMemberAcls(c *gin.Context) {
	var (
		input      = &ListMemberAclsInput{}
		output     = &ListMemberAclsOutput{}
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

	if tagID := c.GetString(cnt.CtxTagID); tagID != "" {
		input.Where = append(input.Where, "tag-id=="+tagID)
	}

	namespace := c.GetString(cnt.CtxNamespace)
	listMemberAclsInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     input.Where,
		Namespace: &namespace,
	}
	listMemberAclsOutput, err := vrm.ListMemberAcls(listMemberAclsInput, c)
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
			zap.String(cnt.Controller, "vrm.ListMemberAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listMemberAclsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.MemberAcls = []MemberAcl{}
	output.Total = int(listMemberAclsOutput.Count)
	for _, data := range listMemberAclsOutput.Data {
		m := MemberAcl{}
		m.ExtractByProto(c, data)
		output.MemberAcls = append(output.MemberAcls, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
