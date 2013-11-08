PRiSM
=====

Work In Progress
----------------
API is unstable.

[![Build Status](https://travis-ci.org/modcloth-labs/prism.png?branch=master)](https://travis-ci.org/modcloth-labs/prism)

RabbitMQ consumer which logs all messages to a local database.

## Usage

Please refer to the `godoc` documentation (http://godoc.org/github.com/modcloth-labs/prism).

Inquiries can be directed to: github+prism@modcloth.com

## Development

Requires SQLite3 header files to be present for testing.

Requires gosqlite package for testing (`go get -v code.google.com/p/gosqlite/sqlite3`)

Test via `go test github.com/modcloth-labs/prism`.