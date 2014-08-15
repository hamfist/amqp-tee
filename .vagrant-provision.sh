#!/bin/bash

export DEBIAN_FRONTEND=noninteractive

set -e
set -x

apt-get update -yq
apt-get install -yq \
  bison \
  sqlite3 \
  libsqlite3-dev \
  libsqlite0\
  rabbitmq-server \
  build-essential \
  git \
  curl \
  mercurial \
  vim-nox

/usr/lib/rabbitmq/bin/rabbitmq-plugins enable rabbitmq_management

if ! docker -v ; then
  curl -s 'https://get.docker.io' | sh
fi

if [[ ! $(go version | grep 1.3) ]] ; then
  curl -s -L 'http://golang.org/dl/go1.3.linux-amd64.tar.gz' | \
    tar xzf - -C /usr/local/
  ln -sfv /usr/local/go/bin/* /usr/local/bin/
fi

mkdir -p /gopath
chown -R vagrant:vagrant /gopath

su - vagrant -c /vagrant/.vagrant-provision-as-vagrant.sh

echo 'Ding!'
