package admin

import (
	auth "VirtualRegistryManagement/authentication"
	authComm "VirtualRegistryManagement/authentication/common"
	"VirtualRegistryManagement/modules/openstack"

	"github.com/Zillaforge/toolkits/tracer"
)

func namespaceIsLegal(namespace string) bool {
	if !openstack.NamespaceIsLegal(namespace) {
		return false
	}
	return true
}
func projectIDIsLegal(projectID string) bool {
	authInput := &authComm.GetProjectInput{
		ID: projectID,
	}
	_, err := auth.Use().GetProject(tracer.StartEntryContext(tracer.EmptyRequestID), authInput)
	if err != nil {
		return false
	}
	return true
}
func userIDIsLegal(userID string) bool {
	authInput := &authComm.GetUserInput{
		ID: userID,
	}
	_, err := auth.Use().GetUser(tracer.StartEntryContext(tracer.EmptyRequestID), authInput)
	if err != nil {
		return false
	}
	return true
}
func membershipIsLegal(projectID, userID string) bool {
	authInput := &authComm.GetMembershipInput{
		UserId:    userID,
		ProjectId: projectID,
	}
	_, err := auth.Use().GetMembership(tracer.StartEntryContext(tracer.EmptyRequestID), authInput)
	if err != nil {
		return false
	}
	return true
}
