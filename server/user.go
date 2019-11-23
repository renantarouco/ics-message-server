package server

// User - Every client connected to the server is a user with relevant information
type User struct {
	Nickname string
	TokenStr string
}

// NewUser - Returns a new user instance
func NewUser(nickname string, tokenStr string) *User {
	return &User{
		Nickname: nickname,
		TokenStr: tokenStr,
	}
}
