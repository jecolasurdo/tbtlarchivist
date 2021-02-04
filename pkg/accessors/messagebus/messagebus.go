package messagebus

// Message is a wrapper around a message body that also provides capability of
// aknowledging receipt of the message.
type Message struct {
	Acknowledger AckNack
	Body         []byte
}

// QueueInfo contains information about the status of a queue.
type QueueInfo struct {
	Messages  int
	Consumers int
}

// A Sender is anything that is capable of transmitting a message to a message
// bus.
type Sender interface {
	Send([]byte) error
	Inspect() (*QueueInfo, error)
}

// A Receiver is anything that is capable of consuming a message from a message
// bus.
type Receiver interface {
	Receive() (*Message, error)
}

// SenderReceiver is anything that is capable of sending and receiving messages
// from a message bus.
type SenderReceiver interface {
	Sender
	Receiver
}

// An AckNack is anything that is able to send positive and negative
// acknowledgements to a message bus.
type AckNack interface {
	Ack() error
	Nack(requeue bool) error
}
