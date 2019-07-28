# Onix Models <img src="../pics/ox.png" width="200" height="200" align="right">

_A model is a collection of __item and link types__ which define a template for storing a domain specific configuration._

For example, in order to store AWS EC2 information, a convention for defining items and their relationship (links) is required. Multiple conventions (models) can be configured to store information about different cloud providers, application platforms, specific infrastructure, software, etc.

This section provides a list of pre-configured models as follows:

| Model | Description |
|---|---|
| [Ansible Inventory](ansible_inventory/readme.md) | Provides the types and rules required to support storage of Ansible inventories. |
| [AWS EC2](aws_ec2/readme.md) | Provides the types and rules required to support storage of AWS EC2 resources. |
| [Kubernetes](k8s/readme.md) | Provides the types and rules to support the recording of Kubernetes namespaces, services, pods, etc. |

([up](../../readme.md))