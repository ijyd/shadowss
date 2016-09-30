package factory

import (
	"fmt"

	"cloud-keeper/pkg/collector"
	"cloud-keeper/pkg/collector/collectorbackend"
)

//Create a storage interface
func Create(c collectorbackend.Config) (collector.Collector, error) {
	switch c.Type {
	case collectorbackend.OperatorsVultr:
		return newVultr(c)
	case collectorbackend.OperatorsDigOC:
		return newDgOcean(c)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", c.Type)
	}
}
