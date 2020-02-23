package utils

import (
	"path"
	"raspi_exporter/common"
	"runtime"
)

func DebugDir() string {
	_, filename, _, _ := runtime.Caller(0)

	return path.Join(path.Dir(path.Dir(filename)), "debug")
}

func UpdateDebugPrerequisite(opts *common.RaspiExpOpts) {
	debugPath := DebugDir()

	opts.UpdateThermalFile(path.Join(debugPath, "temp"))
	opts.UpdateVCGenCmd(path.Join(debugPath, "vcgencmd"))
}
