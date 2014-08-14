#!/bin/bash
set -x
set -e

ln -svf /gopath/src/github.com/modcloth-labs/amqp-tee ~/amqp-tee
ln -svf /vagrant/.vagrant-skel/bashrc /home/vagrant/.bashrc
ln -svf /vagrant/.vagrant-skel/profile /home/vagrant/.profile

source ~/.profile

mkdir -p ~/bin

wget http://guest:guest@localhost:55672/cli/rabbitmqadmin -O ~/bin/rabbitmqadmin
chmod +x ~/bin/rabbitmqadmin
