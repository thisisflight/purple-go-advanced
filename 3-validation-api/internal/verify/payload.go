package verify

type SendRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SendResponse struct {
}

type VerifyResponse struct {
	Token string
}
