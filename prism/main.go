package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modcloth-labs/prism"
	"github.com/streadway/amqp"
)

var (
	logger = log.New(os.Stderr, "[prism] ", log.LstdFlags)

	databaseDriverFlag string
	databaseUriFlag    string
	amqpUriFlag        string
	queueNameFlag      string
)

func init() {
	usageString := `Usage: %s [options]
	Consumes messages from RabbitMQ and writes to database.
`
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageString, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.StringVar(&databaseDriverFlag, "database-driver", "sqlite3", "Database driver to use (possible values: sqlite3, mysql, postgres)")
	flag.StringVar(&databaseUriFlag, "database-uri", "messages.db", "Database uri")
	flag.StringVar(&amqpUriFlag, "amqp-uri", "amqp://guest:guest@localhost:5672/", "AMQP connection URI")
	flag.StringVar(&queueNameFlag, "queue", "default", "Queue to consume from")
	flag.Parse()
}

func main() {
	var (
		amqpConnection *amqp.Connection
		amqpChannel    *amqp.Channel
		deliveryStore  *prism.DeliveryStore
		deliveries     <-chan (amqp.Delivery)
		err            error
	)

	if deliveryStore, err = prism.NewDeliveryStore(databaseDriverFlag, databaseUriFlag); err != nil {
		log.Printf("Could not create message store: %s", err)
		os.Exit(1)
	}

	defer deliveryStore.Close()

	if amqpConnection, err = amqp.Dial(amqpUriFlag); err != nil {
		log.Printf("Could not connect to RabbitMQ: %s", err)
		os.Exit(1)
	}
	defer amqpConnection.Close()

	if amqpChannel, err = amqpConnection.Channel(); err != nil {
		log.Printf("Could not create AMQP channel: %s", err)
		os.Exit(1)
	}

	if deliveries, err = amqpChannel.Consume(queueNameFlag, "", false, false, false, false, nil); err != nil {
		log.Printf("Could not consume from RabbitMQ: %s", err)
		os.Exit(1)
	}

	for delivery := range deliveries {
		log.Printf("Consuming %+v", delivery)

		if err = deliveryStore.Store(&delivery); err != nil {
			log.Printf("Failed to consume: %s", err)
			if err = delivery.Nack(false, true); err != nil {
				log.Printf("Failed to nack: %s", err)
			}
			continue
		}

		if err = delivery.Ack(false); err != nil {
			log.Printf("Failed to ack: %s", err)
		}
	}
}
