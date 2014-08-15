#!/bin/bash
set -x
set -e

mkdir -p ~/bin

ln -svf /gopath/src/github.com/modcloth-labs/amqp-tee ~/
ln -svf /vagrant/.vagrant-skel/bashrc /home/vagrant/.bashrc
ln -svf /vagrant/.vagrant-skel/profile /home/vagrant/.profile

source ~/.profile

wget http://guest:guest@localhost:55672/cli/rabbitmqadmin -O ~/bin/rabbitmqadmin
chmod +x ~/bin/rabbitmqadmin

go get github.com/kr/godep
go get github.com/golang/lint/golint
