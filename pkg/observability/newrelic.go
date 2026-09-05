package observability

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewRelic(appName, licenseKey string) (*newrelic.Application, error) {
	if licenseKey == "" {
		return nil, nil
	}

	return newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(false),
	)
}
