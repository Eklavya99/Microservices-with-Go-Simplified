package main

import (
	"fmt"
	"listener/events"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitMqConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMqConn.Close()
	log.Println("listening for and consuming RabbitMQ messages....")
	consumer, err := events.NewConsumer(rabbitMqConn)
	if err != nil {
		panic(err)
	}
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil{
		log.Println(err)
	}
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ is not yet ready...")
			counts++
		} else {
			connection = conn
			log.Println("Connected to RabbitMQ")
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
