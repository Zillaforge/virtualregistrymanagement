package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1700xxxx: Authentication

	AuthInternalServerErrCode                       = 17000000
	AuthInternalServerErrMsg                        = "internal server error"
	AuthThisAuthenticationTypeIsNotSupportedErrCode = 17000001
	AuthThisAuthenticationTypeIsNotSupportedErrMsg  = "this authentication type is not supported"
	AuthUserHasBeenFrozenErrCode                    = 17000002
	AuthUserHasBeenFrozenErrMsg                     = "the user has been frozen. please contact administrator"
	AuthUserIsNotAnAdministratorErrCode             = 17000003
	AuthUserIsNotAnAdministratorErrMsg              = "the user is not an administrator"
	AuthUserNotFoundErrCode                         = 17000004
	AuthUserNotFoundErrMsg                          = "user not found"
	AuthIncorrectFormatErrCode                      = 17000005
	AuthIncorrectFormatErrMsg                       = "incorrect format of authentication"
	AuthMembershipNotFoundErrCode                   = 17000006
	AuthMembershipNotFoundErrMsg                    = "membership not found"
	AuthProjectNotFoundErrCode                      = 17000007
	AuthProjectNotFoundErrMsg                       = "project not found"
	AuthProjectIsFrozenErrCode                      = 17000008
	AuthProjectIsFrozenErrMsg                       = "project %s is frozen"
	AuthUserIsFrozenErrCode                         = 17000009
	AuthUserIsFrozenErrMsg                          = "user %s is frozen"
	AuthPermissionDeniedErrCode                     = 17000010
	AuthPermissionDeniedErrMsg                      = "permission denied"
	AuthUnmarshalFromCacheErrCode                   = 17000011
	AuthUnmarshalFromCacheErrMsg                    = "unmarshal from cache error"
	AuthPassProjectIDKeyNotFoundErrCode             = 17000012
	AuthPassProjectIDKeyNotFoundErrMsg              = "twcc.%s.paas-project-id key not found"
)

var (
	// 1700xxxx: Authentication

	// 17000000(internal server error)
	AuthInternalServerErr = tkErr.Error(AuthInternalServerErrCode, AuthInternalServerErrMsg)
	// 17000001(this authentication type is not supported)
	AuthThisAuthenticationTypeIsNotSupportedErr = tkErr.Error(AuthThisAuthenticationTypeIsNotSupportedErrCode, AuthThisAuthenticationTypeIsNotSupportedErrMsg)
	// 17000002(the user has been frozen. please contact administrator)
	AuthUserHasBeenFrozenErr = tkErr.Error(AuthUserHasBeenFrozenErrCode, AuthUserHasBeenFrozenErrMsg)
	// 17000003(the user is not an administrator)
	AuthUserIsNotAnAdministratorErr = tkErr.Error(AuthUserIsNotAnAdministratorErrCode, AuthUserIsNotAnAdministratorErrMsg)
	// 17000004(user not found)
	AuthUserNotFoundErr = tkErr.Error(AuthUserNotFoundErrCode, AuthUserNotFoundErrMsg)
	// 17000005(incorrect format of authentication)
	AuthIncorrectFormatErr = tkErr.Error(AuthIncorrectFormatErrCode, AuthIncorrectFormatErrMsg)
	// 17000006(membership not found)
	AuthMembershipNotFoundErr = tkErr.Error(AuthMembershipNotFoundErrCode, AuthMembershipNotFoundErrMsg)
	// 17000007(project not found)
	AuthProjectNotFoundErr = tkErr.Error(AuthProjectNotFoundErrCode, AuthProjectNotFoundErrMsg)
	// 17000008(project %s is frozen)
	AuthProjectIsFrozenErr = tkErr.Error(AuthProjectIsFrozenErrCode, AuthProjectIsFrozenErrMsg)
	// 17000009(user %s is frozen)
	AuthUserIsFrozenErr = tkErr.Error(AuthUserIsFrozenErrCode, AuthUserIsFrozenErrMsg)
	// 17000010(permission denied)
	AuthPermissionDeniedErr = tkErr.Error(AuthPermissionDeniedErrCode, AuthPermissionDeniedErrMsg)
	// 17000011(unmarshal from cache error)
	AuthUnmarshalFromCacheErr = tkErr.Error(AuthUnmarshalFromCacheErrCode, AuthUnmarshalFromCacheErrMsg)
	// 17000012(twcc.%s.paas-project-id key not found)
	AuthPassProjectIDKeyNotFoundErr = tkErr.Error(AuthPassProjectIDKeyNotFoundErrCode, AuthPassProjectIDKeyNotFoundErrMsg)
)
