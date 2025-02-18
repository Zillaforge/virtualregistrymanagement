package iam

import (
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"context"

	"go.uber.org/zap"
	protos "pegasus-cloud.com/aes/pegasusiamclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtil "pegasus-cloud.com/aes/toolkits/utilities"
)

/*
ListProjects returns all of projects and total count from iam server

errors:
- 17000000(internal server error)
*/
func (h *Provider) ListProjects(ctx context.Context, input *authCom.ListProjectsInput) (output *authCom.ListProjectsOutput, err error) {
	ctx, f := tracer.StartWithContext(
		ctx,
		tkUtil.NameOfFunction().String(),
	)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"err":    &err,
	})

	listProjectsInput := &protos.LimitOffset{
		Limit:  input.Limit,
		Offset: input.Offset,
	}
	listProjectsOutput, err := h.poolHandler.ListProjects(listProjectsInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			// TODO: to be defined later
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.ListProjects(...)"),
			zap.String(cnt.RequestID, ctx.Value(cnt.RequestID).(string)),
			zap.Any("input", listProjectsInput),
		).Error(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr)
	}
	output = &authCom.ListProjectsOutput{
		Projects: listProjectsOutput.Data,
		Total:    listProjectsOutput.Count,
	}
	return output, nil
}
