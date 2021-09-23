# Mongo Event Receiver for Pilot Control

An event receiver is a service that processes Syslog events sent by Pilot Control.

One or more receivers can be configured in Pilot Control at the same time.

To implement a receiver simply create a Restful web service with an endpoint that accepts a payload from Pilot Control via POST method.

This [folder](mongo) implements an example of event receiver that uses MongoDB as a data source and adds query endpoints to the Restful API.
