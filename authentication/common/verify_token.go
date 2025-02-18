package common

// VerifyTokenInput defines input structure for VerifyToken
type VerifyTokenInput struct {
	Token string
}

// VerifyTokenOutput defines output structure for VerifyToken
type VerifyTokenOutput struct {
	UserID     string
	Frozen     bool
	Account    string
	SAATUserID string
}
