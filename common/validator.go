package common

import (
	"fmt"
	"os"
	"os/exec"
)

func CheckPrerequisite(opts *RaspiExpOpts) error {
	// 1. Check vcgencmd exists
	if path, err := exec.LookPath("vcgencmd"); err != nil {
		fmt.Println("vcgencmd not found!")
		return err
	} else {
		opts.UpdateVCGenCmd(path)
	}

	// 2. Check Thermal Zone File of raspberry Pi
	if _, err := os.Stat(opts.ThermalZoneFile); os.IsNotExist(err) {
		fmt.Printf("'%s' file doesn't exists!\n", opts.ThermalZoneFile)
		return err
	}

	return nil
}
