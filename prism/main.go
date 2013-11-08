package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "code.google.com/p/gosqlite/sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/modcloth-labs/schema_ensurer"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
)

var (
	migrations = map[string][]string{
		"20131108000000": {`
	  CREATE TABLE IF NOT EXISTS messages(
		uuid char(32),

    content_type character varying(256),
    content_encoding character varying(256),
    delivery_mode smallint,
    priority smallint,
    correlation_id character varying(256),
    reply_to character varying(256),
    expiration character varying(256),
    timestamp timestamp with time zone,
    type character varying(256),
    user_id character varying(256),

    exchange character varying(256),
    routing_key character varying(256),

		body text,

		created_at timestamp without time zone NOT NULL
	  );
	  `,
		},
	}
	insertSql = `
      INSERT INTO messages(
        uuid,
        content_type,
        content_encoding,
        delivery_mode,
        priority,
        correlation_id,
        reply_to,
        expiration,
        timestamp,
        type,
        user_id,
        exchange,
        routing_key,
        body,
        created_at
      ) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  `

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
		db              *sql.DB
		u4              *uuid.UUID
		schemaEnsurer   *sensurer.SchemaEnsurer
		amqpConnection  *amqp.Connection
		amqpChannel     *amqp.Channel
		deliveries      <-chan (amqp.Delivery)
		insertStatement *sql.Stmt
		err             error
	)

	if db, err = sql.Open(databaseDriverFlag, databaseUriFlag); err != nil {
		log.Printf("Could not connect to database: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	schemaEnsurer = sensurer.New(db, migrations, logger)
	if err = schemaEnsurer.EnsureSchema(); err != nil {
		log.Printf("Could not create schema: %s", err)
		os.Exit(1)
	}

	if insertStatement, err = db.Prepare(insertSql); err != nil {
		log.Printf("Could not prepare insert statement: %s", err)
		os.Exit(1)
	}
	defer insertStatement.Close()

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

		if u4, err = uuid.NewV4(); err != nil {
			log.Printf("Could not create GUID: %s", err)
			os.Exit(1)
		}

		if _, err = insertStatement.Exec(
			u4.String(),
			delivery.ContentType,
			delivery.ContentEncoding,
			delivery.DeliveryMode,
			delivery.Priority,
			delivery.CorrelationId,
			delivery.ReplyTo,
			delivery.Expiration,
			delivery.Timestamp,
			delivery.Type,
			delivery.UserId,
			delivery.Exchange,
			delivery.RoutingKey,
			delivery.Body,
			time.Now()); err != nil {
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
