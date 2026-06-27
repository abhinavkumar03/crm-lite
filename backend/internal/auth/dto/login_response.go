package dto

type LoginResponse struct {
	AccessToken string `json:"access_token"`

	ExpiresIn int64 `json:"expires_in"`
}
