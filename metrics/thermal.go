package metrics

import (
	"bytes"
	"raspi_exporter/common"
	"sync"
)

type ThermalAgent struct {
	common.MetricAgent

	cpuMetric *common.Metric
	gpuMetric *common.Metric

	cpuThermal float64
	gpuThermal float64

	cpuMutex sync.Mutex
	gpuMutex sync.Mutex
}

func NewThermalAgent() *ThermalAgent {
	agent := new(ThermalAgent)

	agent.cpuMetric, _ = common.NewMetric("raspi_thermal").AddAttribute("type", "cpu").Build()
	agent.gpuMetric, _ = common.NewMetric("raspi_thermal").AddAttribute("type", "gpu").Build()

	return agent
}

func (agent *ThermalAgent) loadCPUThermal() {
	// TODO Implementing here
}

func (agent *ThermalAgent) loadGPUThermal() {
	// TODO Implementing here
}

func (agent *ThermalAgent) WriteMetrics(buffer *bytes.Buffer) error {
	// TODO Implementing here
	return nil
}

func (agent *ThermalAgent) NeedDaemon() bool {
	// TODO Implementing here
	return true
}

func (agent *ThermalAgent) RunDaemon() {
	// TODO Implementing here
}
