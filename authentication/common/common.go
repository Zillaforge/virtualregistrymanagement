package common

import "context"

type Provider interface {
	VerifySystemAdminToken(ctx context.Context, input *VerifySystemAdminTokenInput) (output *VerifySystemAdminTokenOutput, err error)
	VerifyToken(ctx context.Context, input *VerifyTokenInput) (output *VerifyTokenOutput, err error)
	GetMembership(ctx context.Context, input *GetMembershipInput) (output *GetMembershipOutput, err error)
	GetProject(ctx context.Context, input *GetProjectInput) (output *GetProjectOutput, err error)
	GetUser(ctx context.Context, input *GetUserInput) (output *GetUserOutput, err error)
	ListProjects(ctx context.Context, input *ListProjectsInput) (output *ListProjectsOutput, err error)
	ListMembershipsByProject(ctx context.Context, input *ListMembershipsByProjectInput) (output *ListMembershipsByProjectOutput, err error)
	GetCredential(ctx context.Context, input *GetCredentialInput) (output *GetCredentialOutput, err error)
}
