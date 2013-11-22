# vim:filetype=ruby

Vagrant.configure('2') do |config|
  config.vm.hostname = 'prism'
  config.vm.box = 'precise64'
  config.vm.box_url = 'http://cloud-images.ubuntu.com/vagrant/precise/current/precise-server-cloudimg-amd64-vagrant-disk1.box'

  config.vm.network :private_network, ip: '33.33.33.10', auto_correct: true
  config.vm.network "forwarded_port", guest: 55672, host: 25672
end
