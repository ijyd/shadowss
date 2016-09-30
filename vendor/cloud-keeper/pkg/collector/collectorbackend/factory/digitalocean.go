package factory

import (
	"cloud-keeper/pkg/collector"
	"cloud-keeper/pkg/collector/collectorbackend"
	"cloud-keeper/pkg/collector/digitalocean"
)

func newDgOcean(c collectorbackend.Config) (collector.Collector, error) {
	return digitalocean.NewDigitalOcean(c.APIKey), nil
}
