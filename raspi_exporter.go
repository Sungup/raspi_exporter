package main

import (
	"bytes"
	"fmt"
	"net/http"
	"raspi_exporter/agents"
	"raspi_exporter/common"
)

type MetricHandler struct {
	http.Handler

	metricAgents []common.MetricAgent
}

func newMetricHandler(agentList []common.MetricAgent) *MetricHandler {
	return &MetricHandler{
		metricAgents: agentList,
	}
}

func (h *MetricHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	buffer := bytes.Buffer{}

	for _, agent := range h.metricAgents {
		agent.WriteMetrics(&buffer)
	}

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(buffer.Bytes())
}

func main() {
	// 1. Argument parsing
	opts := common.ArgParse()

	// 2. Add Checking Prerequisite
	if err := common.CheckPrerequisite(opts); err != nil {
		panic(err)
	}

	// 3. Build agents objects
	var metricAgents []common.MetricAgent
	metricAgents = append(metricAgents, agents.NewThermalAgent(opts))

	// 4. Daemonize
	for _, agent := range metricAgents {
		if agent.NeedDaemon() {
			go common.Daemonize(agent, opts)
		}
	}

	// 3. Make and assign handler
	handler := newMetricHandler(metricAgents)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`<h1>Raspi Exporter</h1><a href="/metrics">Metrics</a>`))
	})

	http.Handle("/metrics", handler)

	// 4. Start server
	fmt.Printf("Connect to %s\n", opts.ServerAddr())

	if err := http.ListenAndServe(opts.ListenAddr(), nil); err != nil {
		panic(err)
	}
}
