package agents

import (
	"bytes"
	"io/ioutil"
	"path"
	"raspi_exporter/common"
	"raspi_exporter/utils"
	"testing"
)

const (
	DebugCPUTemp = 41.868
	DebugGPUTemp = 41.9
)

func makeTestAgent() *ThermalAgent {
	opts := new(common.RaspiExpOpts)

	utils.UpdateDebugPrerequisite(opts)

	return NewThermalAgent(opts)
}

func TestLoadCPUThermal(t *testing.T) {
	agent := makeTestAgent()

	// Load test
	if err := agent.loadCPUTemp(); err != nil {
		t.Error(err)
	}

	// Check debug value
	if agent.cpuTemp != DebugCPUTemp {
		t.Errorf("loaded cpu temperature is not same with debug/temp")
	}
}

func TestLoadGPUThermal(t *testing.T) {
	agent := makeTestAgent()

	// Load test
	if err := agent.loadGPUTemp(); err != nil {
		t.Error(err)
	}

	// Check debug value
	if agent.gpuTemp != DebugGPUTemp {
		t.Errorf("loaded gpu temperature is not same with debug/vcgencmd")
	}
}

func TestWriteMetrics(t *testing.T) {
	agent := makeTestAgent()

	var expect string

	// Load Expected Output
	expectFile := path.Join(utils.DebugDir(), "thermal_expect.log")
	if buffer, err := ioutil.ReadFile(expectFile); err != nil {
		t.Errorf("expect file cannot be loaded")
		t.FailNow()
	} else {
		expect = string(buffer)
	}

	// Load CPU Temperature
	if err := agent.loadCPUTemp(); err != nil {
		t.Error(err)
	}

	// Load GPU Temperature
	if err := agent.loadGPUTemp(); err != nil {
		t.Error(err)
	}

	tested := bytes.Buffer{}

	// Output result test.
	agent.WriteMetrics(&tested)
	if tested.String() != expect {
		t.Errorf("tested result is not same with expected result")
	}

	// Re-entrance test.
	tested.Reset()

	agent.WriteMetrics(&tested)
	if tested.String() != expect {
		t.Errorf("re-entered result is not same with expected result")
	}
}

func TestRunDaemon(t *testing.T) {
	agent := makeTestAgent()

	// Re-entrance test.
	_ = agent.RunDaemon()
	_ = agent.RunDaemon()
}
