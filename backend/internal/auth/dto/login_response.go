package dto

type LoginResponse struct {
	AccessToken string `json:"access_token"`

	User UserResponse `json:"user"`
}
