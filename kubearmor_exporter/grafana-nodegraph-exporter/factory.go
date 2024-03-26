package grafananodegraphexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"context"
)

var (
	Type = component.MustNewType("grafana_nodegraph")
)

const (
	LogsStability = component.StabilityLevelBeta
)

// NewFactory creates a factory for the grafana nodegraph exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, LogsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ClientConfig: confighttp.ClientConfig{
			Endpoint: "0.0.0.0:5000",
		},
	}
}

func createLogsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Logs, error) {

	cfg := config.(*Config)
	nodeGraphExporter, _ := newExporter(cfg, set.TelemetrySettings)

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		config,
		nodeGraphExporter.consumeLogs,
		exporterhelper.WithStart(nodeGraphExporter.start),
		exporterhelper.WithShutdown(nodeGraphExporter.stop),
	)

}
