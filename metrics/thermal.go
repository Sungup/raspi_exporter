package metrics

import (
	"bytes"
	"raspi_exporter/common"
)

type ThermalCollector struct {
	common.MetricCollector

	name string
}

func (collector *ThermalCollector) WriteMetrics(buffer *bytes.Buffer) error {
	return nil
}
