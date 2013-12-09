package amqptee

import (
	"database/sql"
	"io/ioutil"
	"log"
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
)

type DeliveryStore struct {
	db              *sql.DB
	insertStatement *sql.Stmt
}

func NewDeliveryStore(databaseDriver string, databaseUri string) (deliveryStore *DeliveryStore, err error) {
	me := &DeliveryStore{}

	if me.db, err = sql.Open(databaseDriver, databaseUri); err != nil {
		return nil, err
	}

	schemaEnsurer := sensurer.New(me.db, migrations, log.New(ioutil.Discard, "", 0))
	if err = schemaEnsurer.EnsureSchema(); err != nil {
		return nil, err
	}

	if me.insertStatement, err = me.db.Prepare(insertSql); err != nil {
		return nil, err
	}

	return me, nil
}

func (me *DeliveryStore) Store(delivery *amqp.Delivery) (err error) {
	var (
		u4 *uuid.UUID
	)

	if u4, err = uuid.NewV4(); err != nil {
		return err
	}

	_, err = me.insertStatement.Exec(
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

func (me *DeliveryStore) Close() {
	me.db.Close()
}
