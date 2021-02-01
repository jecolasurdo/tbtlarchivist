package messagebusiface

// MessageBusMessage is a wrapper around a message body that also provides
// capability of aknowledging receipt of the message.
type MessageBusMessage struct {
	Acknowledger AcknowledgerIface
	Body         []byte
}

// MessageBus is anything that is capable of sending and receiving messages.
type MessageBus interface {
	Send([]byte) error
	Receive() *MessageBusMessage
}

// AcknowledgerIface represents anything that is able to handle positive and
// negative acknowledgements.
type AcknowledgerIface interface {
	Ack() error
	Nack(requeue bool) error
}
