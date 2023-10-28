package responseparams

// UserResponse user info response.
type UserResponse struct {
	UserId string `json:"id,omitempty"`
	Token  string `json:"token"`
}
