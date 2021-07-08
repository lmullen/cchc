package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	dict = "/usr/share/dict/words" // Unix dictionary on MacOS
)

func main() {

	b, err := os.ReadFile(dict)
	if err != nil {
		log.Fatal("Failed to read dictionary", err)
	}

	words := strings.Split(string(b), "\n")

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

	// Only creates a queue if it doesn't already exist
	q, err := ch.QueueDeclare("jobs", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to create a queue: ", err)
	}

	rand.Seed(time.Now().UnixNano())

	// Keep sending random words forever
	for {
		i := rand.Intn(len(words))
		word := words[i]
		sendJob(word, ch, &q)
		sleep := rand.Intn(5)*100 + 500
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

}

func sendJob(word string, channel *amqp.Channel, queue *amqp.Queue) {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(word),
		Timestamp:    time.Now(),
	}

	err := channel.Publish(
		"",
		queue.Name,
		false,
		false,
		msg,
	)
	if err != nil {
		log.Error("Failed to publish message for ", word, ": ", err)
	}
	log.WithFields(log.Fields{
		"word": string(msg.Body),
		"time": msg.Timestamp,
	}).Info("Published message")

}
