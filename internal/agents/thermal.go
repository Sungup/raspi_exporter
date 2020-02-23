package agents

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"raspi_exporter/internal/common"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	VCGenSubCmdGPUTemp = "measure_temp"
)

var (
	VCGenCmdRegExp = regexp.MustCompile("[a-zA-Z='\n\r ]")
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

func command(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	} else {
		return stdout.String(), err
	}
}

func (agent *ThermalAgent) loadCPUTemp() error {
	var buffer []byte
	var err error
	var temp float64

	if buffer, err = ioutil.ReadFile(agent.thermalZone); err != nil {
		return err
	}

	if temp, err = strconv.ParseFloat(strings.TrimSpace(string(buffer)), 32); err != nil {
		return err
	}

	agent.cpuTemp = temp / 1000.0

	return nil
}

func (agent *ThermalAgent) loadGPUTemp() error {
	var buffer string
	var err error
	var temp float64

	if buffer, err = command(agent.vcGenCmd, VCGenSubCmdGPUTemp); err != nil {
		return err
	}

	buffer = VCGenCmdRegExp.ReplaceAllString(buffer, "")

	if temp, err = strconv.ParseFloat(buffer, 16); err != nil {
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
