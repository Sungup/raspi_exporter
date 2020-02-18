package common

import (
	"flag"
	"fmt"
)

type RaspiExpOpts struct {
	Debug bool
	Port  uint
}

const (
	RaspiExporterPort = 9100
)

/*
   Constructor & Argument Parser
*/
func ArgParse() *RaspiExpOpts {
	opts := RaspiExpOpts{}

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

	// 2. Parse arguments
	flag.Parse()

	return &opts
}

func (opts *RaspiExpOpts) ServerAddr() string {
	return fmt.Sprintf("http://localhost:%d", opts.Port)
}

func (opts *RaspiExpOpts) ListenAddr() string {
	return fmt.Sprintf(":%d", opts.Port)
}
