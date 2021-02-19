package messagebustypes

import (
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus/acknowledger"
)

// Message is a wrapper around a message body that also provides capability of
// aknowledging receipt of the message.
type Message struct {
	Acknowledger acknowledger.AckNack
	Body         []byte
}

// QueueInfo contains information about the status of a queue.
type QueueInfo struct {
	Messages  int
	Consumers int
}
