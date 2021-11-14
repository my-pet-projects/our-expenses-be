package domain

// AuthenticationData holds authentication data.
type AuthenticationData struct {
	id   string
	user string
}

// NewAuthenticationData creates a new authentication data.
func NewAuthenticationData(user string) AuthenticationData {
	return AuthenticationData{
		user: user,
	}
}

// ID returns user id.
func (a AuthenticationData) ID() string {
	return a.id
}

// User returns user name.
func (a AuthenticationData) User() string {
	return a.user
}
