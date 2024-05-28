package plugin

import (
	"fmt"
	"strings"
)

type MetricQuery struct {
	MetricID    string         `json:"metricId"`
	Measurement string         `json:"measurement"`
	Filters     []MetricFilter `json:"filters"`
}

type MetricFilter struct {
	Group    string   `json:"group"`
	Operator operator `json:"operator"`
	Value    string   `json:"value"`
}

type operator string

var operatorMap = map[operator]string{
	"=":  "",
	"!=": "!",
}

func (o operator) String() string {
	return operatorMap[o]
}

func (m MetricFilter) String() string {
	return fmt.Sprintf("%s%s:%s", m.Operator.String(), m.Group, m.Value)
}

func (m MetricQuery) FrameName() string {
	metric := m.MetricID
	if len(m.Filters) > 0 {
		cond := strings.Builder{}
		for i, f := range m.Filters {
			if i > 0 {
				cond.WriteString(",")
			}
			cond.WriteString(f.String())
		}
		metric = fmt.Sprintf("%s{%s}", metric, cond.String())
	}
	return fmt.Sprintf("%s(%s)", m.Measurement, metric)
}
