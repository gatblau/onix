# Standard Development VM

The following details a way of spinning up a local Linux development VM under Windows that contains a standard set of tools to help with PilotCtl work (as well as potentially being useful for other tasks!)

It is based on:
- Ubuntu 20.04 LTS (minimal install)
- VirtualBox
- Vagrant

# Requirements
- Windows PC with an OS supported by Vagrant (*It is presumed you will be using Windows 10*)
- Working install of Oracle VirtualBox
  - Download and run the install from the [Virtualbox site](https://www.virtualbox.org/wiki/Downloads)
  - If you are using in a commercial environment, do *not* also add on the VirtualBox Extension Pack, unless you have a valid commercial license from Oracle
- Working install of Vagrant
  - Download and run the install from the [Vagrant site](https://www.vagrantup.com)

# Installation
- Download the latest Vagrantfile from this repo into an empty local directory on your Windows PC
- In the same directory run `vagrant up` - the process may take a few minutes to complete depending on the speed of your PC and internet connection
- When completed
  - run `vagrant ssh` to connect to the VM
  - run `vagrant halt` to shutdown the VM without destroying it
  - run `vagrant up` if you wish to start up the VM again (NB. this will simply start up the VM and will not have to deploy out again)
  - run `vagrant destory` if you wish to completely destroy the VM (NB. this will destroy the VM completely, including any files you may have in the VM filesystem)
  - run `vagrant` for more help if needed

# Additional tool and binaries included
- Git
- Docker
- Docker Compose CLI
- Kubectl CLI
- Tekton CLI
- OpenShift CLI
- Kustomize
- Helm CLI

# Notes
- Ports 80 and 443 are forwarded to the host by default (to 8080 and 8081 respectively). Additional ports can be forwarded by adding additional lines to the Vagrantfile if needed and the VM re-provisioned or re-created.

# Troubleshooting

## Proxy issues blocking download of Vagrant box

If you have issues with your location not being able to pull the Vagrant box down correctly then please download manually and add to Vagrant with the following steps:

1. Browse to the web site hosting the Ubuntu 20.04 box (I.E. https://app.vagrantup.com/generic/boxes/ubuntu2004)
2. Download the *Virtualbox* version of the box to your local hard drive (note that the default filename is a GUID - this is fine)
3. At a command prompt in the directory where your downloaded box is, type "vagrant box add generic/ubuntu2004 `GUID`" (where `GUID` is the name of the file downloaded)