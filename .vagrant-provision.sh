#!/bin/bash

export DEBIAN_FRONTEND=noninteractive

set -e
set -x

apt-get update -yq
apt-get install -yq \
  bison \
  sqlite3 \
  libsqlite0\
  rabbitmq-server \
  build-essential \
  git \
  curl \
  mercurial \
  vim-nox

su - vagrant -c bash <<EOBASH
#!/bin/bash
set -x
set -e

ln -svf /vagrant/.vagrant-skel/bashrc /home/vagrant/.bashrc
ln -svf /vagrant/.vagrant-skel/profile /home/vagrant/.profile

source ~/.profile
if ! which gvm >/dev/null ; then
  set +x
  bash < <(curl -s https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer)
  set -x
fi
source ~/.profile
gvm get
if [ -z "\$(gvm list | grep go1.1.2)" ] ; then
  gvm install go1.1.2
fi

mkdir -p /home/vagrant/gopath/src/github.com/modcloth-labs/
ln -svf /vagrant /home/vagrant/gopath/src/github.com/modcloth-labs/prism

EOBASH
echo 'Ding!'
