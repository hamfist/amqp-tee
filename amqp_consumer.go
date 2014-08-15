package amqptee

import (
	"github.com/streadway/amqp"
)

// AMQPConsumer pulls acssages from an AMQP queue and executes a callback for each
type AMQPConsumer struct {
	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	deliveries     <-chan (amqp.Delivery)
}

// NewAMQPConsumer create a new AMQPConsumer, connect to an AMQP service, and start consuming from the queue
func NewAMQPConsumer(amqpURI string, queueName string) (amqpConsumer *AMQPConsumer, err error) {
	ac := &AMQPConsumer{}

	if ac.amqpConnection, err = amqp.Dial(amqpURI); err != nil {
		return nil, err
	}

	if ac.amqpChannel, err = ac.amqpConnection.Channel(); err != nil {
		return nil, err
	}

	if ac.deliveries, err = ac.amqpChannel.Consume(queueName, "", false, false, false, false, nil); err != nil {
		return nil, err
	}

	return ac, nil
}

// Consume start consuming all messages and execute the callback for each
// if the callback returns an error, this function exits with the same error
// err's if it cannot ack
func (ac *AMQPConsumer) Consume(deliveryHandler func(*amqp.Delivery) (err error)) (err error) {
	for delivery := range ac.deliveries {
		if err = deliveryHandler(&delivery); err != nil {
			return err
		}

		if err = delivery.Ack(false); err != nil {
			return err
		}
	}

	return nil
}

// Close close connections to AMQP service
func (ac *AMQPConsumer) Close() {
	ac.amqpChannel.Close()
	ac.amqpConnection.Close()
}
