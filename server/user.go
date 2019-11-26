package server

// User - Every client connected to the server is a user with relevant information
type User struct {
	Nickname string
	Addr     string
	TokenStr string
}

// NewUser - Returns a new user instance
func NewUser(nickname, addr, tokenStr string) *User {
	return &User{
		Nickname: nickname,
		Addr:     addr,
		TokenStr: tokenStr,
	}
}
