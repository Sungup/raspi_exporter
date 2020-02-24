package agents

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"raspi_exporter/internal/common"
	"strconv"
	"sync"
)

const (
	VCGenSubCmdGPUTemp = "measure_temp"
)

type ThermalAgent struct {
	common.MetricAgent

	vcGenCmd    string // vcgencmd tool path
	thermalZone string // thermal zone file path

	cpuMetric *common.Metric
	gpuMetric *common.Metric

	cpuTemp float64
	gpuTemp float64

	mutex sync.RWMutex
}

func NewThermalAgent(opts *common.RaspiExpOpts) *ThermalAgent {
	agent := new(ThermalAgent)

	agent.vcGenCmd = opts.VCGenCmdPath
	agent.thermalZone = opts.ThermalZoneFile

	agent.cpuMetric, _ = common.NewMetric("raspi_thermal").
		AddAttribute("device", "cpu").
		Build()

	agent.gpuMetric, _ = common.NewMetric("raspi_thermal").
		AddAttribute("device", "gpu").
		Build()

	agent.cpuTemp = 0.0
	agent.gpuTemp = 0.0

	agent.mutex = sync.RWMutex{}

	return agent
}

func command(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return stderr.Bytes(), err
	} else {
		return stdout.Bytes(), err
	}
}

func extractFloat(in []byte) (float64, error) {
	buffer := new(bytes.Buffer)
	buffer.Grow(len(in))

	for i := 0; i < len(in); i++ {
		if in[i] == '.' || ('0' <= in[i] && in[i] <= '9') {
			buffer.WriteByte(in[i])
		} else if buffer.Len() > 0 {
			break
		}
	}

	// Extract 16bit only because Raspi uses very small float value
	return strconv.ParseFloat(buffer.String(), 16)
}

func (agent *ThermalAgent) loadCPUTemp() error {
	var buffer []byte
	var err error
	var temp float64

	if buffer, err = ioutil.ReadFile(agent.thermalZone); err != nil {
		return err
	}

	if temp, err = extractFloat(buffer); err != nil {
		return err
	}

	agent.cpuTemp = temp / 1000.0

	return nil
}

func (agent *ThermalAgent) loadGPUTemp() error {
	var buffer []byte
	var err error
	var temp float64

	if buffer, err = command(agent.vcGenCmd, VCGenSubCmdGPUTemp); err != nil {
		return err
	}

	if temp, err = extractFloat(buffer); err != nil {
		return err
	}

	agent.gpuTemp = temp

	return nil
}

func (agent *ThermalAgent) WriteMetrics(buffer *bytes.Buffer) {
	agent.mutex.RLock()
	defer agent.mutex.RUnlock()

	agent.cpuMetric.WritePromQL(buffer, agent.cpuTemp)
	agent.gpuMetric.WritePromQL(buffer, agent.gpuTemp)
}

func (agent *ThermalAgent) NeedDaemon() bool {
	return true
}

func (agent *ThermalAgent) RunDaemon() error {
	agent.mutex.Lock()
	defer agent.mutex.Unlock()

	// TODO Add error handling
	_ = agent.loadCPUTemp()
	_ = agent.loadGPUTemp()

	return nil
}
