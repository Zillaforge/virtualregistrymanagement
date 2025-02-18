package lbmevents

import (
	auth "VirtualRegistryManagement/authentication"
	authCom "VirtualRegistryManagement/authentication/common"
	"VirtualRegistryManagement/constants"

	"pegasus-cloud.com/aes/toolkits/littlebell"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

func Init() {
	littlebell.InitLittleBell(&littlebell.LittleBellConfig{
		Kind:      constants.Kind,
		Region:    mviper.GetString("littlebell.region"),
		Host:      mviper.GetString("littlebell.host"),
		AccessKey: mviper.GetString("littlebell.access_key"),
		SecretKey: mviper.GetString("littlebell.secret_key"),
		Credential: &littlebell.IAMCredential{
			ProjectID: mviper.GetString("littlebell.credential.project_id"),
			UserID:    mviper.GetString("littlebell.credential.user_id"),
			CallCredentialFunc: func(projectId, userId string) (accessKey string, secretKey string) {
				ctx := tracer.StartEntryContext(tracer.EmptyRequestID)
				getCredentialInput := &authCom.GetCredentialInput{
					UserId:    userId,
					ProjectId: projectId,
				}
				getCredentialOutput, getCredentialErr := auth.Use().GetCredential(ctx, getCredentialInput)
				if getCredentialErr != nil {
					return "", ""
				}
				return getCredentialOutput.AccessKey, getCredentialOutput.SecretKey
			},
		},
		Arn: mviper.GetString("littlebell.arn"),
	})
}
