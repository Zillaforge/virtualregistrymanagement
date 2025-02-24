package constants

const (
	Name           = "VirtualRegistryManagement"
	PascalCaseName = "VirtualRegistryManagement"
	SnakeCaseName  = "virtual_registry_management"
	KebabCaseName  = "virtual-registry-management"
	UpperAbbrName  = "VRM"
	LowerAbbrName  = "vrm"

	Kind                 = PascalCaseName
	Version              = "0.0.6"
	APIPrefix            = "/" + LowerAbbrName + "/api/"
	APIVersion           = "v1"
	GlobalConfigPath     = "etc/ASUS"
	GlobalConfigFilename = KebabCaseName + ".yaml"
	ProductUUIDFilePath  = "/sys/class/dmi/id/product_uuid"

	//Workflow log
	RequestID      = "Request-ID"
	Middleware     = "Middleware"
	Authentication = "Authentication"
	Storage        = "Storage"
	GRPC           = "GRPC"
	Controller     = "Controller"
	Server         = "Server"
	Module         = "Module"
	Cmd            = "Cmd"
	Plugin         = "Plugin"
	EventConsume   = "EventConsume"
	Task           = "Task"

	// Query Keys ...
	QueryToken = "token"

	// Header Keys ...
	HdrHostID                = "Host-ID"
	HdrLocationID            = "Location-ID"
	HdrVersionID             = "Version-ID"
	HdrAuthorization         = "Authorization"
	HdrProjectIDFromKong     = "Project-ID"
	HdrUserIDFromKong        = "User-ID"
	HdrUserRoleFromKong      = "User-Role"
	HdrProjectActiveFromKong = "Project-Active"
	HdrSystemAdminFromKong   = "System-Admin"
	HdrUserAccountFromKong   = "User-Account"
	HdrNamespace             = "X-Namespace"
	HdrSAATUserIDFromKong    = "SAAT-User-ID"

	// Context
	CtxLocationID            = HdrLocationID
	CtxHostID                = HdrHostID
	CtxUserID                = "ctxUserID"
	CtxUserAccount           = "ctxUserAccount"
	CtxProjectID             = "ctxProjectID"
	CtxTenantRole            = "ctxTenantRole"
	CtxOperationName         = "ctxOperationName"
	CtxCreator               = "ctxCreator"
	CtxNamespace             = "ctxNamespace"
	CtxRepositoryID          = "ctxRepositoryID"
	CtxTagID                 = "ctxTagID"
	CtxMemberAclID           = "ctxMemberAclID"
	CtxProjectAclID          = "ctxProjectAclID"
	CtxProjectCountFlag      = "ctxProjectCountFlag" // the soft-limit has been reached if true
	CtxProjectSizeFlag       = "ctxProjectSizeFlag"  // the soft-limit has been reached if true
	CtxProjectUsedCount      = "ctxProjectUsedCount"
	CtxProjectUsedSize       = "ctxProjectUsedSize"
	CtxProjectSoftLimitCount = "ctxProjectSoftLimitCount"
	CtxProjectSoftLimitSize  = "ctxProjectSoftLimitSize"
	CtxVolumeID              = "ctxVolumeID"
	CtxSAATUserID            = "ctxSAATUserID"
	CtxRepositoryProtect     = "ctxRepositoryProtect"
	CtxTagProtect            = "ctxTagProtect"
	CtxExportID              = "ctxExportID"
	CtxVPSMetadata           = "ctxVPSMetadata"
	CtxServerOS              = "ctxServerOS"

	// Params
	ParamProjectID    = "project-id"
	ParamRepositoryID = "repository-id"
	ParamTagID        = "tag-id"
	ParamMemberAclID  = "member-acl-id"
	ParamProjectAclID = "project-acl-id"
	ParamServerID     = "server-id"
	ParamVolumeID     = "volume-id"
	ParamExportID     = "export-id"

	//ReconcileKey ...
	ReconcileKey = "eventpublish"
	SyncKey      = "syncevent"
	AsyncKey     = "asyncevent"
)

type CtxWithExtraVal struct{}
