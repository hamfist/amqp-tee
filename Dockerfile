FROM ubuntu:12.04
MAINTAINER ModCloth

ADD ./.build/amqp-tee /amqp-tee

CMD ["/amqp-tee", "-h"]
