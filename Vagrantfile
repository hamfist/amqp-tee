# vim:filetype=ruby

Vagrant.configure('2') do |config|
  config.vm.hostname = 'amqp-tee'
  config.vm.box = 'hashicorp/precise64'

  config.vm.network :private_network, ip: '33.33.33.10', auto_correct: true
  config.vm.network "forwarded_port", guest: 55672, host: 25672
  config.vm.provision :shell, path: '.vagrant-provision.sh'
  config.vm.define 'amqp-tee' do |host|
  end
  config.vm.synced_folder '.', '/gopath/src/github.com/modcloth-labs/amqp-tee'
end
