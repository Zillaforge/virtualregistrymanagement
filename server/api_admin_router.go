package server

import (
	cnt "VirtualRegistryManagement/constants"
	adminCtl "VirtualRegistryManagement/controllers/api/admin"
	sysCtl "VirtualRegistryManagement/controllers/api/system"
	mid "VirtualRegistryManagement/middlewares/api"

	"github.com/gin-gonic/gin"
	pbac "pegasus-cloud.com/aes/toolkits/pbac/gin"
)

func enableAdminVirtualRegistryManagementRouter(rg *gin.RouterGroup) {
	rg.Use(mid.VerifyAdminAuthentication)
	system := rg.Group("system")
	{
		pbac.GET(system, "versions", sysCtl.GetDetailVersions, cnt.AdminGetDetailVersions.Name, true)
		// api/v1/admin/system/configurations
		pbac.GET(system, "configurations", sysCtl.GetSystemConfigurations, cnt.AdminGetSystemConfigurations.Name, true)
		logs := system.Group("logs")
		{
			pbac.GET(logs, "", sysCtl.ListLogs, cnt.AdminListLogs.Name, true)
			pbac.GET(logs, "download", sysCtl.DownloadLog, cnt.AdminDownloadLog.Name, true)
		}
	}

	rg.Use(mid.VerifyOpstkNamespace)

	pbac.POST(rg, "upload", adminCtl.UploadImage, cnt.AdminUploadImage.Name, false)
	pbac.POST(rg, "import", adminCtl.ImportImage, cnt.AdminImportImage.Name, false)

	pbac.GET(rg, "repositories", adminCtl.ListRepositories, cnt.AdminListRepositories.Name, false)
	repository := rg.Group("repository")
	{
		pbac.POST(repository, "", adminCtl.CreateRepository, cnt.AdminCreateRepository.Name, false)
		repositoryID := repository.Group(":repository-id", mid.VerifyRepository)
		{
			pbac.GET(repositoryID, "", adminCtl.GetRepository, cnt.AdminGetRepository.Name, false)
			pbac.PUT(repositoryID, "", adminCtl.UpdateRepository, cnt.AdminUpdateRepository.Name, false)
			pbac.DELETE(repositoryID, "", adminCtl.DeleteRepository, cnt.AdminDeleteRepository.Name, false)

			protect := repositoryID.Group("protect")
			{
				pbac.POST(protect, "", adminCtl.SetRepositoryProtect, cnt.AdminSetRepositoryProtect.Name, false)
				pbac.DELETE(protect, "", adminCtl.UnsetRepositoryProtect, cnt.AdminUnsetRepositoryProtect.Name, false)
			}
		}
	}

	pbac.GET(rg, "tags", adminCtl.ListTags, cnt.AdminListTags.Name, false)
	tag := rg.Group("tag")
	{
		pbac.POST(tag, "", adminCtl.CreateTag, cnt.AdminCreateTag.Name, false)
		tagID := tag.Group(":tag-id", mid.VerifyTag)
		{
			pbac.GET(tagID, "", adminCtl.GetTag, cnt.AdminGetTag.Name, false)
			pbac.PUT(tagID, "", adminCtl.UpdateTag, cnt.AdminUpdateTag.Name, false)
			pbac.DELETE(tagID, "", adminCtl.DeleteTag, cnt.AdminDeleteTag.Name, false)

			protect := tagID.Group("protect")
			{
				pbac.POST(protect, "", adminCtl.SetTagProtect, cnt.AdminSetTagProtect.Name, false)
				pbac.DELETE(protect, "", adminCtl.UnsetTagProtect, cnt.AdminUnsetTagProtect.Name, false)
			}
		}
		{
			pbac.POST(tagID, "download", adminCtl.DownloadImage, cnt.AdminDownloadImage.Name, false)

			projectAcl := tagID.Group("share")
			{
				pbac.POST(projectAcl, "", adminCtl.CreateProjectAcl, cnt.AdminCreateProjectAcl.Name, false)
				pbac.DELETE(projectAcl, "", adminCtl.DeleteProjectAcl, cnt.AdminDeleteProjectAcl.Name, false)
			}
			{
				pbac.GET(tagID, "memberacls", adminCtl.ListMemberAcls, cnt.AdminListMemberAcls.Name, false)
				pbac.POST(tagID, "memberacl", adminCtl.CreateMemberAclBatch, cnt.AdminCreateMemberAclBatch.Name, false)
			}
		}
	}

	memberAclID := rg.Group("memberacl/:member-acl-id", mid.VerifyMemberAcl)
	{
		pbac.DELETE(memberAclID, "", adminCtl.DeleteMemberAcl, cnt.AdminDeleteMemberAcl.Name, false)
	}

	pbac.GET(rg, "projectacls", adminCtl.ListProjectAcls, cnt.AdminListProjectAcls.Name, false)
	projectAcl := rg.Group("projectacl")
	{
		pbac.POST(projectAcl, "", adminCtl.CreateProjectAcl, cnt.AdminCreateProjectAcl.Name, false)
		projectAclID := projectAcl.Group(":project-acl-id", mid.VerifyProjectAcl)
		{
			pbac.DELETE(projectAclID, "", adminCtl.DeleteProjectAcl, cnt.AdminDeleteProjectAcl.Name, false)
		}
	}

	pbac.GET(rg, "projects", adminCtl.ListProjects, cnt.AdminListProjects.Name, false)
	projectID := rg.Group("project/:project-id", mid.VerifyProject, mid.VerifyProjectResource)
	{
		pbac.GET(projectID, "", adminCtl.GetProjectInfo, cnt.AdminGetProject.Name, false)
		pbac.PUT(projectID, "", adminCtl.UpdateProject, cnt.AdminUpdateProject.Name, false)

		server := projectID.Group("server/:server-id", mid.VerifyOpstkServer)
		{
			pbac.POST(server, "snapshot", adminCtl.CreateSnapshot, cnt.AdminCreateSnapshot.Name, false)
		}
	}

	pbac.GET(rg, "exports", adminCtl.ListExports, cnt.AdminListExports.Name, false)
	exportID := rg.Group("export/:export-id", mid.VerifyExport)
	{
		pbac.GET(exportID, "", adminCtl.GetExport, cnt.AdminGetExport.Name, false)
		pbac.DELETE(exportID, "", adminCtl.DeleteExport, cnt.AdminDeleteExport.Name, false)
	}
}
