package domain

import "errors"

// User holds user data.
type User struct {
	id           string
	username     string
	password     string
	token        string
	refreshToken string
}

// NewUser creates a new user.
func NewUser(id string, username string, password string, token string, refreshToken string) (*User, error) {
	if len(username) == 0 {
		return nil, errors.New("username is empty")
	}
	user := &User{
		id:           id,
		username:     username,
		password:     password,
		token:        token,
		refreshToken: refreshToken,
	}
	return user, nil
}

// ID returns user id.
func (u User) ID() string {
	return u.id
}

// Username returns user name.
func (u User) Username() string {
	return u.username
}

// User returns user password.
func (u User) Password() string {
	return u.password
}

// Token returns user token.
func (u User) Token() string {
	return u.token
}

// RefreshToken returns user refresh token.
func (u User) RefreshToken() string {
	return u.refreshToken
}

// UpdateTokens sets new token and refresh token.
func (u *User) UpdateTokens(token string, refreshToken string) {
	u.token = token
	u.refreshToken = refreshToken
}
