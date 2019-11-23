package server

// Message - Basic message abstraction
type Message struct {
	From string `json:"from"`
	Body string `json:"body"`
}
