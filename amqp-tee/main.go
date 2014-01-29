package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modcloth-labs/amqp-tee"
	"github.com/streadway/amqp"
)

var (
	logger = log.New(os.Stderr, "[amqp-tee] ", log.LstdFlags)

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
	flag.StringVar(&queueNameFlag, "queue", "default", "Queue to consume from (also name of table written to)")
	flag.Parse()
}

func main() {
	var (
		amqpConsumer  *amqptee.AMQPConsumer
		deliveryStore *amqptee.DeliveryStore
		err           error
	)

	if deliveryStore, err = amqptee.NewDeliveryStore(databaseDriverFlag, databaseUriFlag, queueNameFlag); err != nil {
		log.Printf("Could not create message store: %s", err)
		os.Exit(1)
	}

	defer deliveryStore.Close()

	if amqpConsumer, err = amqptee.NewAMQPConsumer(amqpUriFlag, queueNameFlag); err != nil {
		log.Printf("Could not connect to RabbitMQ: %s", err)
		os.Exit(1)
	}
	defer amqpConsumer.Close()

	amqpConsumer.Consume(func(delivery *amqp.Delivery) (err error) {
		log.Printf("Consuming %+v", delivery)

		if err = deliveryStore.Store(delivery); err != nil {
			return err
		}

		return nil
	})

}
