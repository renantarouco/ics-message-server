package server

// User - Every client connected to the server is a user with relevant information
type User struct {
	Nickname string
}

// NewUser - Returns a new user instance
func NewUser(nickname string) *User {
	return &User{
		Nickname: nickname,
	}
}
