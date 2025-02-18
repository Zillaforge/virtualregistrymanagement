package common

import "pegasus-cloud.com/aes/pegasusiamclient/pb"

type ListMembershipsByProjectInput struct {
	ProjectID string
	_         struct{}
}

type ListMembershipsByProjectOutput struct {
	Memberships []*pb.MemberJoin
	Total       int64
	_           struct{}
}
