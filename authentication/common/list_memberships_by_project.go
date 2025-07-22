package common

import "github.com/Zillaforge/pegasusiamclient/pb"

type ListMembershipsByProjectInput struct {
	ProjectID string
	_         struct{}
}

type ListMembershipsByProjectOutput struct {
	Memberships []*pb.MemberJoin
	Total       int64
	_           struct{}
}
