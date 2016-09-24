package factory

import (
	"cloud-keeper/pkg/collector"
	"cloud-keeper/pkg/collector/collectorbackend"
	"cloud-keeper/pkg/collector/vultr"
)

func newVultr(c collectorbackend.Config) (collector.Collector, error) {
	return vultr.NewVultr(c.APIKey), nil
}
