package amqptee

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "code.google.com/p/gosqlite/sqlite3" //Register sqlite3 driver
	_ "github.com/go-sql-driver/mysql"     //Register mysql driver
	_ "github.com/lib/pq"                  //Register postgresql driver
	"github.com/modcloth-labs/schema_ensurer"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
)

var (
	migrationFormats = map[string][]string{
		"20131108000000_%s": {`
	  CREATE TABLE IF NOT EXISTS %s(
		uuid char(36),

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
	insertSqlNormalFormat = `
      INSERT INTO %s(
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
	insertSqlPostgresFormat = `
      INSERT INTO %s(
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
      ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
  `
)

// DeliveryStore represents a backing database to insert the AMQP messages into
type DeliveryStore struct {
	db              *sql.DB
	insertStatement *sql.Stmt
}

// NewDeliveryStore open a connection to the given database and initialize the schema
func NewDeliveryStore(databaseDriver string, databaseURI string, table string) (deliveryStore *DeliveryStore, err error) {
	ds := &DeliveryStore{}

	if ds.db, err = sql.Open(databaseDriver, databaseURI); err != nil {
		return nil, err
	}

	if err = ds.runMigrations(table); err != nil {
		return nil, err
	}

	insertSqlFormat := insertSqlNormalFormat
	if databaseDriver == "postgres" {
		insertSqlFormat = insertSqlPostgresFormat
	}

	if ds.insertStatement, err = ds.db.Prepare(fmt.Sprintf(insertSqlFormat, table)); err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *DeliveryStore) runMigrations(table string) (err error) {
	migrations := map[string][]string{}

	for migrationFormatTag, migrationStatementFormats := range migrationFormats {
		for _, migrationStatementFormat := range migrationStatementFormats {
			migrations[fmt.Sprintf(migrationFormatTag, table)] = append(
				migrations[fmt.Sprintf(migrationFormatTag, table)],
				fmt.Sprintf(migrationStatementFormat, table))
		}
	}

	schemaEnsurer := sensurer.New(ds.db, migrations, log.New(ioutil.Discard, "", 0))
	if err = schemaEnsurer.EnsureSchema(); err != nil {
		return err
	}

	return nil
}

// Store stores the given delivery object in the database
func (ds *DeliveryStore) Store(delivery *amqp.Delivery) (err error) {
	var (
		u4 *uuid.UUID
	)

	if u4, err = uuid.NewV4(); err != nil {
		return err
	}

	_, err = ds.insertStatement.Exec(
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
		time.Now())

	return err
}

// Close closes the connection to the database
func (ds *DeliveryStore) Close() {
	ds.db.Close()
}
