package messagebus

import "github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/messagebustypes"

// A Sender is anything that is capable of transmitting a message to a message
// bus.
type Sender interface {
	Send([]byte) error
	Inspect() (*messagebustypes.QueueInfo, error)
}

// A Receiver is anything that is capable of consuming a message from a message
// bus.
type Receiver interface {
	Receive() (*messagebustypes.Message, error)
}

// SenderReceiver is anything that is capable of sending and receiving messages
// from a message bus.
type SenderReceiver interface {
	Sender
	Receiver
}
