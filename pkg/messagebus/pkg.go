package messagebus

import (
	"fmt"

	"github.com/streadway/amqp"
)

// API is an instance of a message bus. This should to instantiated via the
// Initialize function.
type API struct {
	defaultChannel *amqp.Channel
	queue          *amqp.Queue
}

// Initialize establishes a connection with the underlaying message bus. Once
// a connection is established, the function then verifies or creates a queue
// with the specified name.
func Initialize(queueName string) (*API, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a channel")
	}
	defer ch.Close()

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

	return &API{
		defaultChannel: ch,
		queue:          &q,
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

// Receive retrives a message from the message bus.
func (a *API) Receive() ([]byte, error) {
	panic("not implemented")
}
