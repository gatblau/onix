<img src="../docs/pics/ox.png" width="200" height="200" align="right"/>

# DEEP - Application Landscape Dependency Manager & Deployer

## Background

With the arrival of microservices, applications which were before single monoliths are now made of multiple constituent services.

This explosion in the number of services has increased the complexity in deploying large application landscapes, which have now potentially hundred of microservices.

Consider for example a java application. A jar file can potentially have hundreds of dependent jar files which can form a non-cyclical object graph of dependencies. Solutions like maven or gradle emerged to manage such dependencies and build the required applications.

In the world of microservices, microservice dependency can reach similar levels of complexity as in Java applications as described above. The difference is that instead of building jar files, the challenge now is to deploy and interconnect a network of microservices.

Additionally, not every application out there is a microservice. It is increasingly common to find hybrid landscapes made of legacy (Virtual Machine Deployed) and new (containerised) services. In these cases, deployment of a complete landscape can be required for example, to stand up complete integration testing environments in preparation for release to production.

DEEP (**DE**pendency Manager & d**EP**loyer) is an orchestration microservice, which uses configuration management information in Onix to orchestrate the deployment of entire landscapes on multi-clouds.

The configuration uses an Onix model as described [here](docs/dep_model.md).
