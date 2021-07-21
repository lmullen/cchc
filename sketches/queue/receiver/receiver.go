package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {

	log.Info("Connecting to RabbitMQ")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel: ", err)
	}
	defer ch.Close()

	// Only allow so many messages from the queue at once
	err = ch.Qos(40, 0, true)
	if err != nil {
		log.Fatal("Failed to set prefetch on the message queue: ", err)
	}

	// Only creates a queue if it doesn't already exist
	q, err := ch.QueueDeclare("jobs", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare a queue: ", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register a consumer", err)
	}

	words := make([]string, 100)

	// Repeat endlessly
	for {
		// Get a batch of 100 words
		for i := 0; i < 100; i++ {
			msg := <-msgs
			w := string(msg.Body)
			words[i] = w
			diff := time.Now().Sub(msg.Timestamp)
			log.WithFields(log.Fields{
				"word": w,
				"time": diff,
			}).Info("Received a message")
		}
		res := strings.Join(words, " ")
		fmt.Println("Words: ", res)
		time.Sleep(10 * time.Second) // Simulate work
	}

}
