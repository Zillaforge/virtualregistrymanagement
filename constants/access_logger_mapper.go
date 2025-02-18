package constants

import (
	"fmt"
	"strconv"
)

/*
ServiceID: an int [ 0 - 99 ]
Service.ID 註冊後就無法更改，要是要更改請要通知 UI/UX Team或 Nokia Du(杜岳霖)
註冊 ServiceID請至 gitlab wiki 表格
(<gitlab-url>/aes/pegasus-guide/pegasuscontainerimagebuilder/-/wikis/develop-environment-(docker-compose)#serviceid-table)

actionIDRange: 代表 ActionID 的數量上限
如果 actionIDRange 是 1000, 代表 ActionID: [ 0 - 999 ] + serviceID * 1000
Ex: serviceID = 14, actionIDRange = 1000, ActionID: [ 14000 - 14999 ]
*/
const (
	serviceID     = 24
	actionIDRange = 1000
)

var (
	accessLoggerInfoNameMap map[string]*AccessLoggerInfo = make(map[string]*AccessLoggerInfo)
	accessLoggerInfoIDMap   map[int]*AccessLoggerInfo    = make(map[int]*AccessLoggerInfo)
)

type AccessLoggerInfo struct {
	// ID 固定正數值且不可以重複
	ID int
	// Name 須具備標準格式
	Name string
}

func new(id int, name string) AccessLoggerInfo {
	if id >= actionIDRange || id < 0 {
		id = id % actionIDRange
	}
	id = id + serviceID*actionIDRange
	alInfo := AccessLoggerInfo{id, name}

	if _, exist := accessLoggerInfoNameMap[alInfo.Name]; exist {
		panic(fmt.Sprintf("access logger name is duplicate: %s", alInfo.Name))
	}
	if _, exist := accessLoggerInfoIDMap[alInfo.ID]; exist {
		panic(fmt.Sprintf("access logger id is duplicate: %d", alInfo.ID))
	}
	accessLoggerInfoNameMap[alInfo.Name] = &alInfo
	accessLoggerInfoIDMap[alInfo.ID] = &alInfo
	return alInfo
}

func GetAccessLoggerInfo(name string) *AccessLoggerInfo {
	if v, ok := accessLoggerInfoNameMap[name]; ok {
		return v
	}
	return nil
}

func InsertNewAccessLoggerInfo(actionName string, actionID int) AccessLoggerInfo {
	return new(actionID, actionName)
}

func GetAccessLoggerServiceIDStr() string {
	return strconv.Itoa(serviceID)
}

// AccessLoggerInfo 是用來收集 RESTful API的事件名稱與數字編號的關係。
//
// AccessLoggerInfo.Name須有特定格式來作為不同資源區分。格式會採用如下
//
// <category name>:<operation or action name>
//
// 目前已知的 <category name> 有： user, admin, system分別表示 使用者RESTful
// 管理者RESTful以及系統類型
//
// 特別注意：
// 1). AccessLoggerInfo.ID 註冊後就無法更改，要是要更改請要通知 UI/UX Team或 Nokia Du(杜岳霖)
//

var (
	// Category name : system
	// CategoryName=system通常運作於系統本身或使用者/管理者的間接行為。
	// e.g. SysAsyncPublishProcess AccessLoggerInfo = new(999, "system:SysAsyncPublishProcess")

	// Category name: user
	// 運作於使用者 RESTful API的動作，都需要為 `user:`開頭的前綴
	// e.g. UserVersion AccessLoggerInfo = new(0, "user:Version")
	UserListRepositories     AccessLoggerInfo = new(0, "user:ListRepositories")
	UserCreateRepository     AccessLoggerInfo = new(1, "user:CreateRepository")
	UserGetRepository        AccessLoggerInfo = new(2, "user:GetRepository")
	UserUpdateRepository     AccessLoggerInfo = new(3, "user:UpdateRepository")
	UserDeleteRepository     AccessLoggerInfo = new(4, "user:DeleteRepository")
	UserListTags             AccessLoggerInfo = new(5, "user:ListTags")
	UserCreateTag            AccessLoggerInfo = new(6, "user:CreateTag")
	UserGetTag               AccessLoggerInfo = new(7, "user:GetTag")
	UserUpdateTag            AccessLoggerInfo = new(8, "user:UpdateTag")
	UserDeleteTag            AccessLoggerInfo = new(9, "user:DeleteTag")
	UserListMemberAcls       AccessLoggerInfo = new(10, "user:ListMemberAcls")
	UserCreateMemberAclBatch AccessLoggerInfo = new(11, "user:CreateMemberAclBatch")
	UserDeleteMemberAcl      AccessLoggerInfo = new(12, "user:DeleteMemberAcl")
	UserCreateProjectAcl     AccessLoggerInfo = new(13, "user:CreateProjectAcl")
	UserDeleteProjectAcl     AccessLoggerInfo = new(14, "user:DeleteProjectAcl")
	UserUploadImage          AccessLoggerInfo = new(15, "user:UploadImage")
	UserDownloadImage        AccessLoggerInfo = new(16, "user:DownloadImage")
	UserCreateSnapshot       AccessLoggerInfo = new(17, "user:CreateSnapshot")
	UserVersion              AccessLoggerInfo = new(18, "user:Version")
	UserGetProject           AccessLoggerInfo = new(19, "user:GetProject")
	UserListExports          AccessLoggerInfo = new(20, "user:ListExports")
	UserGetExport            AccessLoggerInfo = new(21, "user:GetExport")
	UserDeleteExport         AccessLoggerInfo = new(22, "user:DeleteExport")

	// Category name: admin
	// 運作於管理者 RESTful API的動作，都需要為 `admin:`開頭的前綴
	// Start ID from 300
	AdminGetDetailVersions       AccessLoggerInfo = new(300, "admin:GetDetailVersions")
	AdminGetSystemConfigurations AccessLoggerInfo = new(301, "admin:GetSystemConfigurations")
	AdminListLogs                AccessLoggerInfo = new(302, "admin:ListLogs")
	AdminDownloadLog             AccessLoggerInfo = new(303, "admin:DownloadLog")
	AdminListRepositories        AccessLoggerInfo = new(304, "admin:ListRepositories")
	AdminCreateRepository        AccessLoggerInfo = new(305, "admin:CreateRepository")
	AdminGetRepository           AccessLoggerInfo = new(306, "admin:GetRepository")
	AdminUpdateRepository        AccessLoggerInfo = new(307, "admin:UpdateRepository")
	AdminDeleteRepository        AccessLoggerInfo = new(308, "admin:DeleteRepository")
	AdminListTags                AccessLoggerInfo = new(309, "admin:ListTags")
	AdminCreateTag               AccessLoggerInfo = new(310, "admin:CreateTag")
	AdminGetTag                  AccessLoggerInfo = new(311, "admin:GetTag")
	AdminUpdateTag               AccessLoggerInfo = new(312, "admin:UpdateTag")
	AdminDeleteTag               AccessLoggerInfo = new(313, "admin:DeleteTag")
	AdminListMemberAcls          AccessLoggerInfo = new(314, "admin:ListMemberAcls")
	AdminListProjectAcls         AccessLoggerInfo = new(315, "admin:ListProjectAcls")
	AdminCreateMemberAclBatch    AccessLoggerInfo = new(316, "admin:CreateMemberAclBatch")
	AdminDeleteMemberAcl         AccessLoggerInfo = new(317, "admin:DeleteMemberAcl")
	AdminCreateProjectAcl        AccessLoggerInfo = new(318, "admin:CreateProjectAcl")
	AdminDeleteProjectAcl        AccessLoggerInfo = new(319, "admin:DeleteProjectAcl")
	AdminUploadImage             AccessLoggerInfo = new(320, "admin:UploadImage")
	AdminDownloadImage           AccessLoggerInfo = new(321, "admin:DownloadImage")
	AdminImportImage             AccessLoggerInfo = new(322, "admin:ImportImage")
	AdminCreateSnapshot          AccessLoggerInfo = new(323, "admin:CreateSnapshot")
	AdminGetProject              AccessLoggerInfo = new(324, "admin:GetProject")
	AdminListProjects            AccessLoggerInfo = new(325, "admin:ListProjects")
	AdminUpdateProject           AccessLoggerInfo = new(326, "admin:UpdateProject")
	AdminSetRepositoryProtect    AccessLoggerInfo = new(327, "admin:SetRepositoryProtect")
	AdminSetTagProtect           AccessLoggerInfo = new(328, "admin:SetTagProtect")
	AdminUnsetRepositoryProtect  AccessLoggerInfo = new(329, "admin:UnsetRepositoryProtect")
	AdminUnsetTagProtect         AccessLoggerInfo = new(330, "admin:UnsetTagProtect")
	AdminListExports             AccessLoggerInfo = new(331, "admin:ListExports")
	AdminGetExport               AccessLoggerInfo = new(332, "admin:GetExport")
	AdminDeleteExport            AccessLoggerInfo = new(333, "admin:DeleteExport")
)
