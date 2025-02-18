package common

type GetCredentialInput struct {
	UserId    string
	ProjectId string
	_         struct{}
}

type GetCredentialOutput struct {
	AccessKey string
	SecretKey string
	_         struct{}
}
