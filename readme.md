# Onix 

Onix is a lightweight configuration management database (CMDB) designed to support [Infrastructure as a Code](https://en.wikipedia.org/wiki/Infrastructure_as_code) provisioning to ultimately provide a [single source of thruth](https://en.wikipedia.org/wiki/Single_source_of_truth) across [multi-cloud](https://en.wikipedia.org/wiki/Multicloud) environments.
<img src="docs/pics/ox.png" width="250" height="250" align="right">

The key features are:
- accessible via a [RESTful Web API](./docs/wapi.md), [Ansible](https://www.ansible.com/) and [Terraform](https://www.terraform.io/) compoments ([Connectors](./connectors/readme.md)). A user interface is coming soon. :smiley:
- flexible customisation via [Models](./models/readme.md)  
- fully containerised, runs natively on Openshift and Kubernetes
- automatic population to avoid data misalignments, via Feeders is planned for a near future release.

The following topics give more insight into Onix:

- [Architecture](./docs/architecture.md)
- [Getting started](./docs/getting_started.md)
- [Web API](./docs/wapi.md)
- [Connectors](./connectors/readme.md)
- [Models](./models/readme.md)

## License Terms

This software is licensed under the [Apache License - Version 2.0, January 2004](http://www.apache.org/licenses/).

Copyright (c) 2018-2019 by [gatblau.org](http://gatblau.org).

#### Contributor Notice

Contributors to this project, hereby assign copyright in their code to the 
project, to be licensed under the same terms as the rest of the code.