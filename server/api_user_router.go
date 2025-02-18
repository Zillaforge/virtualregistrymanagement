package server

import (
	cnt "VirtualRegistryManagement/constants"
	sysCtl "VirtualRegistryManagement/controllers/api/system"
	userCtl "VirtualRegistryManagement/controllers/api/user"
	mid "VirtualRegistryManagement/middlewares/api"

	"github.com/gin-gonic/gin"
	pbac "pegasus-cloud.com/aes/toolkits/pbac/gin"
)

func enableUserVirtualRegistryManagementRouter(rg *gin.RouterGroup) {
	pbac.GET(rg, "version", sysCtl.GetPlainTextVersion, cnt.UserVersion.Name, false)

	rg.Use(mid.VerifyOpstkNamespace)
	rg.Use(mid.VerifyUserToken)

	projectID := rg.Group("project/:project-id", mid.VerifyMembership, mid.VerifyProjectResource)
	{
		pbac.GET(projectID, "", userCtl.GetProjectInfo, cnt.UserGetProject.Name, false)
		pbac.POST(projectID, "upload", userCtl.UploadImage, cnt.UserUploadImage.Name, false)

		server := projectID.Group("server/:server-id", mid.VerifyOpstkServer)
		{
			pbac.POST(server, "snapshot", userCtl.CreateSnapshot, cnt.UserCreateSnapshot.Name, false)
		}

		pbac.GET(projectID, "repositories", userCtl.ListRepositories, cnt.UserListRepositories.Name, false)
		repository := projectID.Group("repository")
		{
			pbac.POST(repository, "", userCtl.CreateRepository, cnt.UserCreateRepository.Name, false)
			repositoryID := repository.Group(":repository-id", mid.VerifyRepositoryHasInProject)
			{
				pbac.GET(repositoryID, "", userCtl.GetRepository, cnt.UserGetRepository.Name, false)
				pbac.PUT(repositoryID, "", userCtl.UpdateRepository, cnt.UserUpdateRepository.Name, false)
				pbac.DELETE(repositoryID, "", userCtl.DeleteRepository, cnt.UserDeleteRepository.Name, false)

				pbac.GET(repositoryID, "tags", userCtl.ListTags, cnt.UserListTags.Name, false)
				pbac.POST(repositoryID, "tag", userCtl.CreateTag, cnt.UserCreateTag.Name, false)

				pbac.POST(repositoryID, "memberacl", userCtl.CreateMemberAclBatch, cnt.UserCreateMemberAclBatch.Name, false)
			}
		}
		pbac.GET(projectID, "tags", userCtl.ListTags, cnt.UserListTags.Name, false)
		tagID := projectID.Group("tag/:tag-id")
		{
			aclTagID := tagID.Group("", mid.VerifyTagHasInACL)
			{
				pbac.GET(aclTagID, "", userCtl.GetTag, cnt.UserGetTag.Name, false)
			}

			projectTagID := tagID.Group("", mid.VerifyTagHasInProject)
			{
				pbac.PUT(projectTagID, "", userCtl.UpdateTag, cnt.UserUpdateTag.Name, false)
				pbac.DELETE(projectTagID, "", userCtl.DeleteTag, cnt.UserDeleteTag.Name, false)

				projectAcl := projectTagID.Group("share")
				{
					pbac.POST(projectAcl, "", userCtl.CreateProjectAcl, cnt.UserCreateProjectAcl.Name, false)
					pbac.DELETE(projectAcl, "", userCtl.DeleteProjectAcl, cnt.UserDeleteProjectAcl.Name, false)
				}
				pbac.POST(projectTagID, "download", userCtl.DownloadImage, cnt.UserDownloadImage.Name, false)

				{
					pbac.GET(projectTagID, "memberacls", userCtl.ListMemberAcls, cnt.UserListMemberAcls.Name, false)
					pbac.POST(projectTagID, "memberacl", userCtl.CreateMemberAclBatch, cnt.UserCreateMemberAclBatch.Name, false)
				}
			}
		}

		memberAclID := projectID.Group("memberacl/:member-acl-id", mid.VerifyMemberAclHasInProject)
		{
			pbac.DELETE(memberAclID, "", userCtl.DeleteMemberAcl, cnt.UserDeleteMemberAcl.Name, false)
		}

		pbac.GET(projectID, "exports", userCtl.ListExports, cnt.UserListExports.Name, false)
		exportID := projectID.Group("export/:export-id", mid.VerifyExportHasInProject)
		{
			pbac.GET(exportID, "", userCtl.GetExport, cnt.UserGetExport.Name, false)
			pbac.DELETE(exportID, "", userCtl.DeleteExport, cnt.UserDeleteExport.Name, false)
		}
	}
}
