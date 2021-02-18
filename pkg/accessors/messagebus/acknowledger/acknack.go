package acknowledger

// An AckNack is anything that is able to send positive and negative
// acknowledgements to a message bus.
type AckNack interface {
	Ack() error
	Nack(requeue bool) error
}
