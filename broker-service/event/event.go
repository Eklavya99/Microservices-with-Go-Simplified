package events

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error{
	return ch.ExchangeDeclare(
		"log_topic", // Name
		"topic", //type
		true, //durable?
		false, //auto-delete?
		false, //internal?
		false, // no-wait?
		nil, //arguements
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error){
	return ch.QueueDeclare(
		"", //name
		false, //durable?
		false, //delete when not in use
		true, //exclusive?
		false, //no wait??
		nil, // arguements?
	)
} 