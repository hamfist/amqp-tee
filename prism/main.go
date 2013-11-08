package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "code.google.com/p/gosqlite/sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/modcloth-labs/schema_ensurer"
	"github.com/streadway/amqp"
)

var (
	migrations = map[string][]string{
		"20131108000000": {`
	  CREATE TABLE IF NOT EXISTS messages(
		id serial PRIMARY KEY,

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

		body text,

		created_at timestamp without time zone NOT NULL
	  );
	  `,
		},
	}
	logger = log.New(os.Stderr, "[prism] ", log.LstdFlags)

	databaseDriverFlag string
	databaseUriFlag    string
	amqpUriFlag        string
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
	flag.Parse()
}

func main() {
	var (
		db             *sql.DB
		schemaEnsurer  *sensurer.SchemaEnsurer
		amqpConnection *amqp.Connection
		amqpChannel    *amqp.Channel
		err            error
	)

	if db, err = sql.Open(databaseDriverFlag, databaseUriFlag); err != nil {
		log.Printf("Could not connect to database: %s", err)
		os.Exit(1)
	}

	schemaEnsurer = sensurer.New(db, migrations, logger)
	if err = schemaEnsurer.EnsureSchema(); err != nil {
		log.Printf("Could not create schema: %s", err)
		os.Exit(1)
	}

	if amqpConnection, err = amqp.Dial(amqpUriFlag); err != nil {
		log.Printf("Could not connect to RabbitMQ: %s", err)
		os.Exit(1)
	}
	defer amqpConnection.Close()

	if amqpChannel, err = amqpConnection.Channel(); err != nil {
		log.Printf("Could not create AMQP channel: %s", err)
		os.Exit(1)
	}

  _ = amqpChannel
}
