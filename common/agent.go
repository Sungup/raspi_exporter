package common

import (
	"bytes"
	"time"
)

type MetricAgent interface {
	WriteMetrics(buffer *bytes.Buffer)
	NeedDaemon() bool
	RunDaemon() error
}

func Daemonize(metricAgent MetricAgent, opts *RaspiExpOpts) {
	duration := time.Duration(opts.Refresh) * time.Second

	for true {
		// TODO Handling error
		_ = metricAgent.RunDaemon()
		time.Sleep(duration)
	}
}
