package common

import "pegasus-cloud.com/aes/pegasusiamclient/pb"

type ListProjectsInput struct {
	Limit  int32
	Offset int32
	_      struct{}
}

type ListProjectsOutput struct {
	Projects []*pb.ProjectInfo
	Total    int64
	_        struct{}
}
