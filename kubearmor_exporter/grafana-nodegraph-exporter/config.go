package grafananodegraphexporter

import (
	"fmt"
	"go.opentelemetry.io/collector/config/confighttp"
	"net/url"
)

type Config struct {
	confighttp.ClientConfig `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.

	RelayEndpoint string `mapstructure:"relayendpoint,omitempty"`
	LogFilter     string `mapstructure:"logfilter,omitempty"`
}

func Validate(c *Config) error {
	if _, err := url.Parse(c.Endpoint); c.Endpoint == "" || err != nil {

		return fmt.Errorf("\"endpoint\" must be a valid URL")
	}

	if _, err := url.Parse(c.RelayEndpoint); c.RelayEndpoint == "" || err != nil {

		return fmt.Errorf("\"KubeArmor relay endpoint\" must be a valid URL")
	}

	if c.LogFilter == "" {
		return fmt.Errorf("Log filter should have a valid value")
	}

	return nil
}
