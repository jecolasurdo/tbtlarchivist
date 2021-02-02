package amqpadapter

import "github.com/streadway/amqp"

// Acknowledger provides capabilities for ack, nack, or reject a message that
// has been received.
type Acknowledger struct {
	channel *amqp.Channel
	tag     uint64
}

// NewAcknowledger returns an instance of an acknowledger associated with a
// specific message channel and message tag.
func NewAcknowledger(channel *amqp.Channel, tag uint64) *Acknowledger {
	return &Acknowledger{
		channel: channel,
		tag:     tag,
	}
}

// Ack acknowledges that a message has been received.
func (a *Acknowledger) Ack() error {
	return a.channel.Ack(a.tag, false)
}

// Nack negatively acknowledges that a message has been received.
func (a *Acknowledger) Nack(requeue bool) error {
	return a.channel.Nack(a.tag, false, requeue)
}
