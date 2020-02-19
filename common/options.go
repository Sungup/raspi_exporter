package common

import (
	"flag"
	"fmt"
)

type RaspiExpOpts struct {
	Debug   bool
	Port    uint
	Refresh uint

	VCGenCmdPath    string
	ThermalZoneFile string
}

const (
	RaspiExporterPort    = 9100
	DefaultRefresh       = 2
	ThermalZonePath      = "/sys/class/temp/thermal_zone0/temp"
	DebugThermalZonePath = "debug/temp"
)

/*
   Constructor & Argument Parser
*/
func ArgParse() *RaspiExpOpts {
	opts := RaspiExpOpts{
		VCGenCmdPath:    "",
		ThermalZoneFile: "",
	}

	// 1. Assign arguments
	flag.BoolVar(
		&opts.Debug,
		"debug",
		false,
		"Run on the debug mode",
	)

	flag.UintVar(
		&opts.Port,
		"port",
		RaspiExporterPort,
		"Exporter web port",
	)

	flag.UintVar(
		&opts.Refresh,
		"refresh",
		DefaultRefresh,
		"Exporter refresh duration",
	)

	// 2. Parse arguments
	flag.Parse()

	if opts.Debug {
		opts.ThermalZoneFile = DebugThermalZonePath
	} else {
		opts.ThermalZoneFile = ThermalZonePath
	}

	return &opts
}

func (opts *RaspiExpOpts) ServerAddr() string {
	return fmt.Sprintf("http://localhost:%d", opts.Port)
}

func (opts *RaspiExpOpts) ListenAddr() string {
	return fmt.Sprintf(":%d", opts.Port)
}

func (opts *RaspiExpOpts) UpdateVCGenCmd(path string) {
	opts.VCGenCmdPath = path
}
