package common

type GetUserInput struct {
	ID        string
	Cacheable bool
}

type GetUserOutput struct {
	ID          string
	DisplayName string
	Account     string
	Email       string
	Frozen      bool
	Extra       map[string]interface{}
}
