package utils

import (
	"path"
	"raspi_exporter/internal/common"
	"runtime"
)

func DebugDir() string {
	_, filename, _, _ := runtime.Caller(0)

	return path.Join(path.Dir(path.Dir(path.Dir(filename))), "test")
}

func UpdateDebugPrerequisite(opts *common.RaspiExpOpts) {
	debugPath := DebugDir()

	opts.UpdateThermalFile(path.Join(debugPath, "temp"))
	opts.UpdateVCGenCmd(path.Join(debugPath, "vcgencmd"))
}
