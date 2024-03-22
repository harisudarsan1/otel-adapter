package grafananodegraphexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"context"
	"time"
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
		BackOffConfig: configretry.NewDefaultBackOffConfig(),
		QueueSettings: exporterhelper.NewDefaultQueueSettings(),
		DefaultLabelsEnabled: map[string]bool{
			"exporter": true,
			"job":      true,
			"instance": true,
			"level":    true,
		},
	}
}

type Config struct {
	confighttp.ClientConfig      `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings `mapstructure:"sending_queue"`
	configretry.BackOffConfig    `mapstructure:"retry_on_failure"`

	DefaultLabelsEnabled map[string]bool `mapstructure:"default_labels_enabled"`
}

func createLogsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Logs, error) {

}
