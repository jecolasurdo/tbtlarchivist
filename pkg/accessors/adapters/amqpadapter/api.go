package amqpadapter

import (
	"context"
	"fmt"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/streadway/amqp"
)

// API is an instance of a message bus. This should to instantiated via the
// Initialize function.
type API struct {
	defaultChannel *amqp.Channel
	queue          *amqp.Queue
	inboundMsgs    <-chan amqp.Delivery
}

// Initialize establishes a connection with the underlaying message bus. Once
// a connection is established, the function then verifies or creates a queue
// with the specified name. The API then immediately begins receiving messages
// from the queue.
func Initialize(ctx context.Context, queueName string, prefetchCount int) (*API, error) {
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

	err = ch.Qos(prefetchCount, 0, false)
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

// Send transmits a message to the message bus.
func (a *API) Send(msg []byte) error {
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

// Receive retrieves a message from the message bus. This method does not
// block.  If no message is available the method will return nil.
func (a *API) Receive() *messagebus.Message {
	select {
	case msg := <-a.inboundMsgs:
		acknowledger := NewAcknowledger(a.defaultChannel, msg.DeliveryTag)
		return &messagebus.Message{
			Body:         msg.Body,
			Acknowledger: acknowledger,
		}
	default:
		return nil
	}
}
