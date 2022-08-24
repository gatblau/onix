# Open Telemetry Collector for Pilot

This section contains a telemetry collector for Pilot using OpenTelemetry.

It implements the following  pipeline:

[hostmetrics] --> [fileexporter]
[syslogs]     --> [fileexporter]

- [Entry point for testing](collector/collector_test.go)
- [Sample configuration](telem.yaml)