package stater

// Message types
const (
	DoneMessage MessageType = "done"
)

// MessageType is a way to quickly classify the nature of the message.
//
// The type can be something standard like the constants above, or custom
type MessageType string

// Message is a message sent from some task
type Message struct {
	Type    MessageType
	Task    *Task
	Message interface{}
}

// Messager is something every task has access to that the task uses to send
// messages.
// By default the code sends a message when the task is paused, errored, or finished
type Messager struct {
	messageStream chan *Message
}

// NewMessager creates an initialized messager
func NewMessager() *Messager {
	return &Messager{
		messageStream: make(chan *Message),
	}
}

// GetMessageStream gets the channel of messages.
//
// Note that some thread should be actively listening ot this channel otherwise
// messages will be dropped.
// Messages are sent non-blocking
func (m *Messager) GetMessageStream() chan *Message {
	return m.messageStream
}

// SendMessage sends a message to the message stream without blocking.
// This means if nobody is listening to the message it will not be received
func (m *Messager) SendMessage(message *Message) {
	select {
	case m.messageStream <- message:
	default:
	}
}
