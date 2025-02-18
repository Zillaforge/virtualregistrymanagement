package common

type VerifySystemAdminTokenInput struct {
	Token string
	_     struct{}
}

type VerifySystemAdminTokenOutput struct {
	UserID  string
	Account string
	_       struct{}
}
