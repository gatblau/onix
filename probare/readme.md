<img src="static/img/probe.svg" width="150" height="150" align="right">

# Probare - An application configuration prototype

Probare is a generic web application that is designed to demonstrate changes made to its configuration by Onix.

It also helps modelling how different applications load configuration data and allow to test and show how Onix can be used to manage any application configuration.

Its web user interface displays a terminal showing application events in real time, a textbox containing any configuration file(s) content and a list of environment variables the application can see.

In order to represent applications whose configuration wants to be managed by Onix, it implements different types of configuration loading techniques as follows:

## Configuration Loading Types

Probare can load configuration via:

1. **Environment variables**: this is the most common way to load configuration information to containers.

2. **Configuration file(s)**: some applications require the configuration information to be stored in one or more files. Files can also contain environment variables for merging.

3. **HTTP endpoint**: for applications having sensitive information, one way to make it more difficult to gain aunothorised access to application configuration is by not exposing it through environment variables or storing it in files. An http endpoint can allow the application to load the configuration from the posted http payload.

## Configuration Loading Triggers

Probare also implement different techniques to trigger the reloading of configuration:

1. **Application Restart**: configuration is loaded upon application start. If the application does not implement any logic to trigger configuration reloads then this is the option to use.

2. **POSIX Signalling**: the application listen to a POSIX signal to determine if changed configurations should be reloaded. Probare reloads configuration if it receives a SIGHUP (terminal hang up) signal.

3. **HTTP endpoint (no payload)**: calling an http endpoint with no payload triggers the application to reload its configuration.

4. **HTTP endpoint (configuration payload)**: calling an http endpoint but posting the configuration as payload.