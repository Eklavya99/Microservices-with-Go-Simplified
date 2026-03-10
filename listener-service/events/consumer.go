package events

import (
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct{
	conn *amqp.Connection
	queueName string
}

type Payload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn : conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)	
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	que, err := declareRandomQueue(ch)
	if err != nil{
		return err
	}

	for _, s := range topics{
		ch.QueueBind(
			que.Name,
			s,
			"log_topic",
			false,
			nil,
		)
	}

	messages, err := ch.Consume(que.Name, "", true, false, false, false, nil)
	if err != nil{
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range messages{
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()
	fmt.Printf("Waiting for messages [Exchange, Queue] [log_topic, %s]\n", que.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload){
	switch payload.Name{
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "atuh":

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}		

	}
}

func logEvent(entry Payload) error{
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil{
		return err
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}