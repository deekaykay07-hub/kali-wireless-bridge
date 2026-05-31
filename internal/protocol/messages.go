package protocol

// MessageType defines the type of message exchanged between bridge and TUI
type MessageType string

const (
	TypeCapabilities MessageType = "capabilities"
	TypeCommand      MessageType = "command"
	TypeOutput       MessageType = "output"
	TypeResult       MessageType = "result"
	TypeError        MessageType = "error"
)

// Message is the base message format
type Message struct {
	Type    MessageType `json:"type"`
	Payload any         `json:"payload"`
}

// CommandPayload is sent from TUI to bridge to execute something
type CommandPayload struct {
	ID      string   `json:"id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Privileged bool   `json:"privileged,omitempty"`
}

// OutputPayload streams real-time output back to the TUI
type OutputPayload struct {
	ID     string `json:"id"`
	Stream string `json:"stream"` // stdout or stderr
	Data   string `json:"data"`
}

// ResultPayload is the final result of a command
type ResultPayload struct {
	ID       string `json:"id"`
	ExitCode int    `json:"exit_code"`
	Success  bool   `json:"success"`
}
