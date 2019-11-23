package server

const (
	// CommandMessage - Command to broadcast a message
	CommandMessage = "message"
	// CommandSetNickname - Command to change user's nickname
	CommandSetNickname = "setnick"
	// CommandSwitchRoom - Command to switch user's room
	CommandSwitchRoom = "switchroom"
	// CommandCreateRoom - Command to create a new room
	CommandCreateRoom = "createroom"
	// CommandListUsers - Command to list users within same room
	CommandListUsers = "listusers"
	// CommandListRooms - Command that lists all available rooms
	CommandListRooms = "listrooms"
	// CommandExit - Command that gracefully disconnects a client
	CommandExit = "exit"
)

// Command - Structure representing clint commands
type Command struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}
