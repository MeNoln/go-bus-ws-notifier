package busclient

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gopkg.in/olahol/melody.v1"
)

func ProcessMessages(m *melody.Melody) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				log.Fatal(err.Error())
			}
		}
	}()

	conn, err := amqp.Dial(viper.GetString("amqpbus.host"))
	throwOnError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	throwOnError(err)
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		viper.GetString("amqpbus.queueName"),
		true,
		false,
		false,
		false,
		nil,
	)
	throwOnError(err)

	err = ch.Qos(
		1,
		0,
		false,
	)
	throwOnError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	throwOnError(err)

	lock := make(chan bool)

	go func() {
		for msg := range msgs {
			processMessage(&msg, m)
		}
	}()

	log.Info("Consuming messages")
	<-lock
}

func processMessage(msg *amqp.Delivery, m *melody.Melody) {
	var body CurrencyCreatedEvent
	json.Unmarshal(msg.Body, &body)

	sockMsg := fmt.Sprintf("[%s] Currency with name %s was created.", body.MsgId, body.Title)

	log.WithField("body", sockMsg).Info("Received message from bus: ")
	m.Broadcast([]byte(sockMsg))

	msg.Ack(false)
}

func throwOnError(e error) {
	if e != nil {
		panic(e)
	}
}
