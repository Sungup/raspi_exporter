package metrics

import (
	"bytes"
	"io/ioutil"
	"path"
	"raspi_exporter/common"
	"runtime"
	"testing"
)

const (
	DebugCPUTemp = 41.868
	DebugGPUTemp = 41.9
)

func debugDir() string {
	_, filename, _, _ := runtime.Caller(0)

	return path.Join(path.Dir(path.Dir(filename)), "debug")
}

func makeTestAgent() *ThermalAgent {
	debugPath := debugDir()

	opts := common.RaspiExpOpts{}
	opts.UpdateThermalFile(path.Join(debugPath, "temp"))
	opts.UpdateVCGenCmd(path.Join(debugPath, "vcgencmd"))

	return NewThermalAgent(&opts)
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
	if buffer, err := ioutil.ReadFile(path.Join(debugDir(), "thermal_expect.log")); err != nil {
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
