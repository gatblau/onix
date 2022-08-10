/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package artisanfileexporter

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "file"
	// The stability level of the exporter.
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for OTLP exporter.
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporterAndStabilityLevel(createTracesExporter, stability),
		component.WithMetricsExporterAndStabilityLevel(createMetricsExporter, stability),
		component.WithLogsExporterAndStabilityLevel(createLogsExporter, stability))
}

func createDefaultConfig() config.Exporter {
	return &fileexporter.Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
	}
}

func createTracesExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.TracesExporter, error) {
	fe := exporters.GetOrAdd(cfg, func() component.Component {
		return &fileExporter{path: cfg.(*fileexporter.Config).Path}
	})
	return exporterhelper.NewTracesExporter(
		cfg,
		set,
		fe.Unwrap().(*fileExporter).ConsumeTraces,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithShutdown(fe.Shutdown),
	)
}

func createMetricsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.MetricsExporter, error) {
	fe := exporters.GetOrAdd(cfg, func() component.Component {
		return &fileExporter{path: cfg.(*fileexporter.Config).Path}
	})
	return exporterhelper.NewMetricsExporter(
		cfg,
		set,
		fe.Unwrap().(*fileExporter).ConsumeMetrics,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithShutdown(fe.Shutdown),
	)
}

func createLogsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.LogsExporter, error) {
	fe := exporters.GetOrAdd(cfg, func() component.Component {
		return &fileExporter{path: cfg.(*fileexporter.Config).Path}
	})
	return exporterhelper.NewLogsExporter(
		cfg,
		set,
		fe.Unwrap().(*fileExporter).ConsumeLogs,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithShutdown(fe.Shutdown),
	)
}

// This is the map of already created File exporters for particular configurations.
// We maintain this map because the Factory is asked trace and metric receivers separately
// when it gets CreateTracesReceiver() and CreateMetricsReceiver() but they must not
// create separate objects, they must use one Receiver object per configuration.
var exporters = NewSharedComponents()
