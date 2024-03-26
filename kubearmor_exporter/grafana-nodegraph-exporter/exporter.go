package grafananodegraphexporter

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
)

type nodeGraphExporter struct {
	config          *Config
	settings        component.TelemetrySettings
	wg              sync.WaitGroup
	kubearmorClient *Feeder
	cancel          context.CancelFunc
}

func newExporter(config *Config, settings component.TelemetrySettings) (*nodeGraphExporter, error) {
	settings.Logger.Info("using the grafana exporter")
	logClient, _ := NewClient(config.RelayEndpoint, config.LogFilter)
	return &nodeGraphExporter{
		config:          config,
		settings:        settings,
		kubearmorClient: logClient,
	}, nil

}

func (n *nodeGraphExporter) consumeLogs(ctx context.Context, ld plog.Logs) error {
	return nil
}

func (nge *nodeGraphExporter) start(_ context.Context, host component.Host) (err error) {

	ctx, cancel := context.WithCancel(context.Background())
	nge.cancel = cancel
	nge.wg = sync.WaitGroup{}

	nge.wg.Add(1)
	go nge.nodeGraph(ctx)
	nge.wg.Add(1)
	go nge.startServer(ctx)

	return nil
}

func (nge *nodeGraphExporter) stop(context.Context) (err error) {

	nge.cancel()

	if err := nge.kubearmorClient.DestroyClient(); err != nil {
		return fmt.Errorf("Failed to destroy the kubearmor relay gRPC client (%s)\n", err.Error())
	}
	nge.wg.Wait()
	return nil
}
