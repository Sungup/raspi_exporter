package common

import (
	"os"
	"os/exec"
)

const (
	ThermalZoneFile = "/sys/class/temp/thermal_zone0/temp"
)

func CheckPrerequisite(opts *RaspiExpOpts) error {
	// 1. Check vcgencmd exists
	if path, err := exec.LookPath("vcgencmd"); err != nil {
		return err
	} else {
		opts.UpdateVCGenCmd(path)
	}

	// 2. Check Thermal Zone File of raspberry Pi
	if _, err := os.Stat(ThermalZoneFile); os.IsNotExist(err) {
		return err
	} else {
		opts.UpdateThermalFile(ThermalZoneFile)
	}

	return nil
}
