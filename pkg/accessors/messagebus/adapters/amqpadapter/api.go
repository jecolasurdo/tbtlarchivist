package amqpadapter

import (
	"context"
	"fmt"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/streadway/amqp"
)

// Direction denotes if the message queue allows sending and/or receiving.
type Direction int

const (
	// DirectionReceiveOnly callers will only receive messages from this queue.
	DirectionReceiveOnly Direction = 1

	// DirectionSendOnly callers will only send messages to this queue.
	DirectionSendOnly Direction = 2
)

// API is an instance of a message bus. This should to instantiated via the
// Initialize function.
type API struct {
	defaultChannel *amqp.Channel
	queue          *amqp.Queue
	inboundMsgs    <-chan amqp.Delivery
	direction      Direction
}

// Initialize establishes a connection with the underlaying message bus. Once a
// connection is established, the function then verifies or creates a queue
// with the specified name. If the supplied direction is DirectionReceiveOnly
// then the API will immediately begin receiving messages from the queue.
func Initialize(ctx context.Context, queueName string, direction Direction) (*API, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a channel")
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue")
	}

	if direction == DirectionSendOnly {
		return &API{
			defaultChannel: ch,
			queue:          &q,
			inboundMsgs:    nil,
		}, nil
	}

	err = ch.Qos(5, 0, false)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		return nil, err
	}

	// monitor the channel.Notify methods and a parent context.Done for
	// close requests, and call ch.Close if warranted
	// defer ch.Close()

	return &API{
		defaultChannel: ch,
		queue:          &q,
		inboundMsgs:    msgs,
	}, nil
}

// Send transmits a message to the message bus. This method will panic if the
// message bus was initialized as receive-only.
func (a *API) Send(msg []byte) error {
	if a.direction == DirectionReceiveOnly {
		panic("Cannot send on a receive-only connection.")
	}

	err := a.defaultChannel.Publish(
		"",           // exchange
		a.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		err = fmt.Errorf("Failed to publish a message")
	}
	return err
}

// Inspect returns information about the number of messages and consumers
// associated with the queue.
func (a *API) Inspect() (*messagebus.QueueInfo, error) {
	info, err := a.defaultChannel.QueueInspect(a.queue.Name)
	if err != nil {
		return nil, err
	}
	return &messagebus.QueueInfo{
		Messages:  info.Messages,
		Consumers: info.Consumers,
	}, nil
}

// Receive retrieves a message from the message bus. This method does not
// block.  If no message is available the method will return nil. This method
// will panic if the message bus was initialized as send-only.
func (a *API) Receive() (*messagebus.Message, error) {
	if a.direction == DirectionSendOnly {
		panic("Cannot receive from a send-only connection.")
	}

	select {
	case msg, open := <-a.inboundMsgs:
		if !open {
			return nil, fmt.Errorf("message bus is closed")
		}
		acknowledger := NewAcknowledger(a.defaultChannel, msg.DeliveryTag)
		return &messagebus.Message{
			Body:         msg.Body,
			Acknowledger: acknowledger,
		}, nil
	default:
		return nil, nil
	}
}

var _ messagebus.SenderReceiver = (*API)(nil)
