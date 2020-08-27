---
id: onix_arc
title: Architecture
---
import useBaseUrl from '@docusaurus/useBaseUrl';

Onix is software that proactively records and controls IT systems configuration. 

For modularity and simplicity of application to specific problems, it is designed as a set of lightweight,
modular and independent components. 

These components can be grouped into different areas, depending on their function, as follows:

1. **Core**: core components provide a data model, a procedure for creating, updating, deleting configuration data 
and issuing change notifications. They are the [object-relational database](https://en.wikipedia.org/wiki/Object-relational_database), 
the [Web API](https://en.wikipedia.org/wiki/Web_API) and the [message broker](https://en.wikipedia.org/wiki/Message_broker) 
using a [publish-subscribe pattern](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern).

2. **Web API extensions**: these components modularly add extra endpoints to the Web API to suit particular use cases. 

3. **Clients**: allow different technologies to update and retrieve configuration data connecting to the Web API.

4. **Configuration Managers**: manage the configuration of specific IT infrastructure or applications.

They can be seen in the following picture:

<img alt="Onix Architecture" src={useBaseUrl('img/ox_solution.png')} />
