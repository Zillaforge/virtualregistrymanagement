package admin

import (
	"VirtualRegistryManagement/authentication"
	authComm "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type CreateMemberAclBatchInput struct {
	TagID  []string `json:"tagID" binding:"required"`
	UserID []string `json:"userID" binding:"required"`

	_ struct{}
}

type CreateMemberAclBatchOutput struct {
	MemberAcls []MemberAcl `json:"memberAcls"`
	_          struct{}
}

func CreateMemberAclBatch(c *gin.Context) {
	var (
		input      = &CreateMemberAclBatchInput{}
		output     = &CreateMemberAclBatchOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusCreated
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindJSON(input); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindJSON(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if tagID := c.GetString(cnt.CtxTagID); len(input.TagID) == 0 && tagID != "" {
		input.TagID = append(input.TagID, tagID)
	}

	projectID := c.GetString(cnt.CtxProjectID)
	membership := map[string]bool{}
	{
		authInput := &authComm.ListMembershipsByProjectInput{
			ProjectID: projectID,
		}
		authOutput, err := authentication.Use().ListMembershipsByProject(c, authInput)
		if err != nil {
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cnt.AuthMembershipNotFoundErr.Code():
					statusCode = http.StatusForbidden
					err = tkErr.New(cnt.MidPermissionDeniedErr)
					utility.ResponseWithType(c, statusCode, err)
					return
				}
			}
			zap.L().With(
				zap.String(cnt.Controller, "authentication.Use().ListMembershipsByProject(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", authInput),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.AdminAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}

		for _, member := range authOutput.Memberships {
			membership[member.UserID] = true
		}
	}

	tags := map[string]bool{}
	{
		namespace := c.GetString(cnt.CtxNamespace)
		listRegistriesInput := &pb.ListRegistriesInput{
			Limit:  -1,
			Offset: 0,
			Flag: &pb.RegistryFlag{
				ProjectID: &projectID,

				BelongProject: true,
			},
			Namespace: &namespace,
		}
		listRegistriesOutput, err := vrm.ListRegistries(listRegistriesInput, c)
		if err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "vrm.ListRegistries(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", listRegistriesInput),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.AdminAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}

		for _, registry := range listRegistriesOutput.Data {
			tags[registry.TagID] = true
		}
	}

	createMemberAclBatchInput := &pb.MemberAclBatchInfo{}
	for _, tagID := range input.TagID {
		if _, exist := tags[tagID]; !exist {
			statusCode = http.StatusForbidden
			err = tkErr.New(cnt.AdminAPITagNotBelongRepositoryErr, tagID)
			utility.ResponseWithType(c, statusCode, err)
			return
		}

		for _, userID := range input.UserID {
			if _, exist := membership[userID]; !exist {
				statusCode = http.StatusForbidden
				err = tkErr.New(cnt.AdminAPIUserNotBelongProjectErr, userID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}

			createMemberAclBatchInput.Data = append(createMemberAclBatchInput.Data, &pb.MemberAclInfo{
				TagID:  tagID,
				UserID: userID,
			})
		}
	}

	createMemberAclBatchOutput, err := vrm.CreateMemberAclBatch(createMemberAclBatchInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateMemberAclBatch(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createMemberAclBatchInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.MemberAcls = []MemberAcl{}
	for _, data := range createMemberAclBatchOutput.Data {
		m := MemberAcl{}
		m.ExtractByProto(c, data)
		output.MemberAcls = append(output.MemberAcls, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
