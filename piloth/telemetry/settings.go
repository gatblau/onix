/*
  Onix Config Manager - Pilot
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package telemetry

import (
	"github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
	"os"
)

// NewSettings returns new settings for the collector with default values.
func NewSettings(configPaths []string, version string, loggingOpts []zap.Option) (*service.CollectorSettings, error) {
	// configure receivers
	receiverMap, err := component.MakeReceiverFactoryMap(
		hostmetricsreceiver.NewFactory(),
	)
	if err != nil {
		return nil, err
	}
	// configure processors
	processorMap, err := component.MakeProcessorFactoryMap(
		// add labels to metrics
		resourceattributetransposerprocessor.NewFactory(),
		// add a unique (host.name) to the metric resource(s), allowing users to filter between multiple systems
		resourcedetectionprocessor.NewFactory(),
		// aggregates incoming metrics into a batch, releasing them if a certain time has passed or if a certain number
		// of entries have been aggregated
		batchprocessor.NewFactory(),
	)
	if err != nil {
		return nil, err
	}
	// configure exporters
	exporterMap, err := component.MakeExporterFactoryMap(
		prometheusexporter.NewFactory(),
	)
	if err != nil {
		return nil, err
	}
	buildInfo := component.BuildInfo{
		Command:     os.Args[0],
		Description: "piloth open-telemetry collector for host metrics",
		Version:     version,
	}
	// reads the configuration from a file
	fileP := fileprovider.New()
	configProviderSettings := service.ConfigProviderSettings{
		Locations:     configPaths,
		MapProviders:  map[string]confmap.Provider{fileP.Scheme(): fileP},
		MapConverters: []confmap.Converter{expandconverter.New()},
	}
	provider, err := service.NewConfigProvider(configProviderSettings)
	if err != nil {
		return nil, err
	}
	return &service.CollectorSettings{
		Factories: component.Factories{
			Receivers:  receiverMap,
			Processors: processorMap,
			Exporters:  exporterMap,
		},
		BuildInfo:               buildInfo,
		LoggingOptions:          loggingOpts,
		ConfigProvider:          provider,
		DisableGracefulShutdown: true,
	}, nil
}
