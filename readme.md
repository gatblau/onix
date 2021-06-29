# Onix Configuration Manager 

<<<<<<< HEAD
Onix is a lightweight configuration manager system  designed to support [Infrastructure as a Code](https://en.wikipedia.org/wiki/Infrastructure_as_code) provisioning to ultimately provide a [single source of truth](https://en.wikipedia.org/wiki/Single_source_of_truth) across [multi-cloud](https://en.wikipedia.org/wiki/Multicloud) environments.
<img src="docs/pics/ox.png" width="200" height="200" align="right">

## High Level Overview

For a **high level overview** of the application and break down of its current and planned components take a look at the [overview](docs/overview.md) section.

For a summary of the main **development themes** and features take a look at this [section](docs/dev_themes.md).

## Contributions & feedback

Contributions are welcome, see [here](CONTRIBUTING.MD) for more information.

## Key Features
=======
Onix is a lightweight configuration management database (CMDB) designed to support [Infrastructure as a Code](https://en.wikipedia.org/wiki/Infrastructure_as_code) provisioning to ultimately provide a [single source of truth](https://en.wikipedia.org/wiki/Single_source_of_truth) across [multi-cloud](https://en.wikipedia.org/wiki/Multicloud) environments.
<img src="docs/pics/ox.png" width="300" height="300" align="right">
>>>>>>> upstream/master

The key features are:
- accessible via [RESTful Web API](./docs/wapi.md), [Ansible](https://www.ansible.com/) and [Terraform](https://www.terraform.io/) components ([Connectors](./connectors/readme.md)) and a [Web Console](wc/readme.md).
- flexible customisation via [Models](docs/models/readme.md)  
- fully containerised, runs natively on Openshift and Kubernetes
- automatic population to avoid data misalignments, e.g. Kubernetes

The following topics give more insight into Onix:

- [Architecture](./docs/architecture.md)
- [Getting started](./docs/getting_started.md)
- [Database deployment](./docs/db_deploy.md)
- [Web API](./docs/wapi.md)
- [API Extensions](./docs/extensions.md)
- [Role Based Access](./docs/rbac.md)
  - [OpenId Connect](./docs/oidc.md)
- [Connectors](./connectors/readme.md)
- [Models](docs/models/readme.md)
- [Web Console](wc/readme.md)

## License Terms

This software is licensed under the [Apache License - Version 2.0, January 2004](http://www.apache.org/licenses/).

Copyright (c) 2018-2019 by [gatblau.org](http://gatblau.org).

## Contributor Notice

Contributors to this project, hereby assign copyright in their code to the project, to be licensed under the same terms as the rest of the code.

## Acknowledgements

<a href="https://www.jetbrains.com/?from=onix"><img src="docs/pics/jetbrains.png" width="150" align="right"></img></a>

<<<<<<< HEAD
Many thanks to [JetBrains](https://www.jetbrains.com/?from=onix) for providing an open source license for this project.
=======
Contributors to this project, hereby assign copyright in their code to the 
project, to be licensed under the same terms as the rest of the code.
>>>>>>> upstream/master
