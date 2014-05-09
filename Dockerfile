FROM ubuntu:12.04
MAINTAINER ModCloth

RUN apt-get update -y

RUN apt-get install -y \
  libsqlite3-dev

ADD ./.build/amqp-tee /amqp-tee

CMD ["/amqp-tee", "-h"]
