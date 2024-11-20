# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  # Test environments
  {
    'ubuntu' => 'ubuntu/focal64',
    'fedora' => 'fedora/35-cloud-base',
    'debian' => 'debian/bullseye64'
  }.each do |name, box|
    config.vm.define name do |machine|
      machine.vm.box = box
      machine.vm.hostname = "diceware-#{name}"

      # Resources
      machine.vm.provider "virtualbox" do |vb|
        vb.memory = "2048"
        vb.cpus = 2
      end

      # Shared folder
      machine.vm.synced_folder ".", "/vagrant", type: "virtualbox"

      # Provisioning
      machine.vm.provision "shell", inline: <<-SHELL
        # Install required packages
        if command -v apt-get &> /dev/null; then
          apt-get update
          apt-get install -y podman make git
        elif command -v dnf &> /dev/null; then
          dnf -y install podman make git
        fi

        # Test build and run
        cd /vagrant
        make build
        # Add automated tests here
      SHELL
    end
  end
end
