package prism

import (
	"github.com/streadway/amqp"
)

type AMQPConsumer struct {
	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	deliveries     <-chan (amqp.Delivery)
}

func NewAMQPConsumer(amqpUri string, queueName string) (amqpConsumer *AMQPConsumer, err error) {
	me := &AMQPConsumer{}

	if me.amqpConnection, err = amqp.Dial(amqpUri); err != nil {
		return nil, err
	}

	if me.amqpChannel, err = me.amqpConnection.Channel(); err != nil {
		return nil, err
	}

	if me.deliveries, err = me.amqpChannel.Consume(queueName, "", false, false, false, false, nil); err != nil {
		return nil, err
	}

	return me, nil
}

func (me *AMQPConsumer) Consume(deliveryHandler func(*amqp.Delivery) (err error)) (err error) {
	for delivery := range me.deliveries {
		if err = deliveryHandler(&delivery); err != nil {
			return err
		}

		if err = delivery.Ack(false); err != nil {
			return err
		}
	}

	return nil
}

func (me *AMQPConsumer) Close() {
	me.amqpChannel.Close()
	me.amqpConnection.Close()
}
