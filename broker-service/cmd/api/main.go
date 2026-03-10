package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"math"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct{
	Rabbit *amqp.Connection
}

func main() {
	rabbitMqConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMqConn.Close()	
	app := Config{
		Rabbit: rabbitMqConn,
	}
	log.Printf("Starting broker service on port %s\n", webPort)

	//define Http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort), // ":80"
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
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
