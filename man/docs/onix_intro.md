---
id: onix_intro
title: Introduction
---

Onix is a set of cloud native services that work together to achieve three key goals:
record, report and control IT infrastructure and applications on hybrid clouds.

## Recording any configuration change automatically

Recording is an important aspect of compliance. Modern IT landscapes are usually complex
as they are typically comprised of hundreds or even thousands of granular services.
Meeting compliance requirements means been able to accurately, reliably, comprehensively and 
securely record changes in the configuration of IT systems. With an ever growing population of IT 
services and platforms, keeping their configuration information up to date is becoming more difficult.

Onix is a reactive configuration manager, built as native cloud software, which proactively records any configuration change 
automatically without the intervention of an operator. This ensures that configuration 
information is always up to date. 

Not only the latest configuration is recorded but 
also every change that ever happen to it, providing a comprehensive audit trail that can be 
used for compliance, non-repudiation, troubleshooting, management and security purposes.

## Report any change of the configuration anywhere

Modern IT landscapes are typically deployed across multiple cloud providers both public and private.
Every cloud provider has their own set of tools to monitor and report on their own resources
but with the increasing number of options out there, it becomes almost impossible to visualise everything that happens 
on a hybrid landscape in real time and from a single location. 

Onix acts as an aggregator for configuration information across providers, platforms and applications; allowing to 
report comprehensively on not only the current but any past inventories of the landscape at any point in time.

## Control IT systems

Recording configuration information is only one aspect of configuration management.
Should a change in a IT system configuration be recorded after the system was changed?
What if a change in the configuration could trigger an actual change in the IT system.
In other words, what if the configuration system can act as a proactive controller?

Onix strives to do both. In certain cases it is better to be lazy and record the change
after the system has changed. There are other instances however, where driving a change in the 
actual system from a change in the configuration is more advantageous.

This is the concept used by Kubernetes. Onix leverages this concept keeping any Kubernetes 
configuration centrally in a change history and then issuing marshalling those changes to the 
Kubernetes API.

Another example of control is for instance when Onix pushes changes in application configuration
to services operating on the cloud. This forces applications to alter their behaviour from a 
well curated and trusted configuration source.