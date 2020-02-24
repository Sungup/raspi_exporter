package common

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	valStart = `="`
	valEnd   = `"`
)

type Metric struct {
	name       string
	attributes map[string]string

	buffer []byte
}

func NewMetric(name string) *Metric {
	metric := new(Metric)

	metric.name = name
	metric.attributes = make(map[string]string)

	return metric
}

func (metric *Metric) AddAttribute(key string, value string) *Metric {
	metric.attributes[key] = value

	return metric
}

func (metric *Metric) Build() (*Metric, error) {
	if metric.name == "" {
		return metric, errors.New("metric name is empty")
	}

	buffer := bytes.NewBufferString(metric.name)

	if len(metric.attributes) > 0 {
		buffer.WriteByte('{')

		var splitter = ""

		for key, val := range metric.attributes {
			buffer.WriteString(splitter)
			buffer.WriteString(key)
			buffer.WriteString(valStart)
			buffer.WriteString(val)
			buffer.WriteString(valEnd)

			if len(splitter) == 0 {
				splitter = ", "
			}
		}

		buffer.WriteString("} ")
	}

	metric.buffer = buffer.Bytes()

	return metric, nil
}

func (metric *Metric) WritePromQL(buffer *bytes.Buffer, value float64) {
	buffer.Write(metric.buffer)

	_, _ = fmt.Fprintf(buffer, "%.3f\n", value)
}
