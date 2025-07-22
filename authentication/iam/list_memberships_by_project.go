package iam

import (
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"context"

	"go.uber.org/zap"
	protos "github.com/Zillaforge/pegasusiamclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtil "github.com/Zillaforge/toolkits/utilities"
)

/*
ListMembershipsByProject returns all of projects and total count from iam server

errors:
- 17000000(internal server error)
*/
func (h *Provider) ListMembershipsByProject(ctx context.Context, input *authCom.ListMembershipsByProjectInput) (output *authCom.ListMembershipsByProjectOutput, err error) {
	ctx, f := tracer.StartWithContext(
		ctx,
		tkUtil.NameOfFunction().String(),
	)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"err":    &err,
	})

	listMembershipsByProjectInput := &protos.ListMembershipByProjectInput{
		ProjectID: input.ProjectID,
		Data: &protos.LimitOffset{
			Limit:  -1,
			Offset: 0,
		},
	}
	listMembershipsByProjectOutput, err := h.poolHandler.ListMembershipsByProject(listMembershipsByProjectInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			// TODO: to be defined later
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.ListMembershipsByProject(...)"),
			zap.String(cnt.RequestID, ctx.Value(cnt.RequestID).(string)),
			zap.Any("input", listMembershipsByProjectInput),
		).Error(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr)
	}
	output = &authCom.ListMembershipsByProjectOutput{
		Memberships: listMembershipsByProjectOutput.Data,
		Total:       listMembershipsByProjectOutput.Count,
	}
	return output, nil
}
