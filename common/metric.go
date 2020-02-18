package common

import (
	"bytes"
	"strconv"
)

const (
	valStart = `="`
	valEnd   = `"`
)

type Metric struct {
	attributes map[string]string
	value      float64
}

func (metric *Metric) PromQL() []byte {
	buffer := &bytes.Buffer{}
	buffer.Grow(256)

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

		buffer.WriteByte('}')
	}

	buffer.WriteString(strconv.FormatFloat(metric.value, 'f', 3, 32))

	return buffer.Bytes()
}

type MetricCollector interface {
	WriteMetrics(buffer *bytes.Buffer) error
}
