amqp-tee
========

Work In Progress
----------------
API is unstable.

[![Build Status](https://travis-ci.org/modcloth-labs/amqp-tee.png?branch=master)](https://travis-ci.org/modcloth-labs/amqp-tee)

AMQP consumer which will consume all messages from a given queue and insert
them into a configurable database with a static schema (you can see the schema
in `delivery_store.go`). Currently it only stores the properties and body of
the message, but not any headers.

Run `amqp-tee -h` for a listing of flags.

## Usage

Please refer to the `godoc` documentation (http://godoc.org/github.com/modcloth-labs/amqp-tee).

Inquiries can be directed to: github+amqptee@modcloth.com

## Development

Requires SQLite3 header files to be present for testing.

Requires gosqlite package for testing (`go get -v code.google.com/p/gosqlite/sqlite3`)

Test via `go test github.com/modcloth-labs/amqptee`.
