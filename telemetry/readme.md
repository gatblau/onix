# Open Telemetry Collector for Pilot

This section contains an experimental collector for Pilot using OpenTelemetry.

Currently not integrated, implements the following  pipeline:

[hostmetrics] --> [prometheus]

- [Entry point for testing](opentelem_test.go)
- [Sample configuration](telem.yaml)